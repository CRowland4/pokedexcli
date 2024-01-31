package pokeapi

import (
	"net/http"
	"fmt"
	"errors"
	"encoding/json"
	"io/ioutil"
	"github.com/CRowland4/pokedexcli/internal/pokecache"
	"sync"
)

const locationCount = 20

// Struct to read in the response from the LocationAreas endpoint of the PokéAPI
type locationArea struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func LocationGetter() (getLocations func(string) ([locationCount]string, error)) {
	currentLocationID := 1

	getLocations = func(command string) (locations [locationCount]string, err error) {
		if command == "mapb" && currentLocationID == 1 {
			return locations, errors.New("No previous locations!")
		}
		cacheAllLocationsIfNotCached(currentLocationID, command)

		switch command {
		case "mapb":
			for i := range locations {
				currentLocationID--
				locations[i], _ = pokecache.Get(currentLocationID)
			}
		case "map":
			for i := range locations {
				locations[i], _ = pokecache.Get(currentLocationID)
				currentLocationID++
			}
		default:
			return locations, errors.New(fmt.Sprintf("Location getter command not recognized: %s", command))
		}

		return locations, nil
	}

	return getLocations
}

func cacheAllLocationsIfNotCached(currentLocation int, command string) {
	var wg sync.WaitGroup

	switch command {
	case "mapb":
		for i := 0; i < locationCount; i++ {
			currentLocation--
			go cacheLocationIfNotCached(currentLocation, &wg)
		}
	case "map":
		for i := 0; i < locationCount; i++ {
			go cacheLocationIfNotCached(currentLocation, &wg)
			currentLocation++
		}
	default:
		fmt.Println("Command not recognized for caching map locations:", command)
	}

	wg.Wait()
	return
}

func cacheLocationIfNotCached(id int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if pokecache.IsCached(id) {
		return
	}

	apiResponse := getPokeAPILocation(id)
	pokecache.Add(id, apiResponse)
	return
}

func getPokeAPILocation(id int) (locationName string) {
	address := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", id)
	response, errResponse := http.Get(address)
	if errResponse != nil {
		return fmt.Sprintf("Error retrieving location id %d: %w", id, errResponse)
	}
	defer response.Body.Close()

	body, errBody := ioutil.ReadAll(response.Body)
	if errBody != nil {
		return fmt.Sprintf("Unable to read body of location API response: %w", errBody)
	}

	var locationResponse locationArea
	errUnmarshal := json.Unmarshal(body, &locationResponse)
	if errUnmarshal != nil {
		return fmt.Sprintf("Unable to unmarshal location API response for ID %d: %w", id, errUnmarshal)
	}

	return locationResponse.Name
}