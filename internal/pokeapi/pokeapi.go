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

const LocationCount = 20

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

func LocationCacher() (cacheLocations func(*pokecache.Cache, string) ([LocationCount]string, error)) {
	currentLocationID := 1
	
	cacheLocations = func(cache *pokecache.Cache, command string) (locations [LocationCount]string, err error) {
		if command == "mapb" && currentLocationID <= LocationCount + 1 {
			return locations, errors.New("No previous locations!")
		} else if command == "mapb" {
			currentLocationID -= (2 * LocationCount)
		}

		cacheAllLocationsIfNotCached(cache, currentLocationID, command)
		locations = getCachedLocations(*cache, currentLocationID, command)
		currentLocationID += LocationCount
		return locations, nil
	}

	return cacheLocations
}

func getCachedLocations(cache pokecache.Cache, locationID int, command string) (locations [LocationCount]string) {
	for i := 0; i < LocationCount; i++ {
		entry, _ := cache.Get(locationID)
		locations[i] = entry.LocationName
		locationID++
	}
	return locations
}

func cacheAllLocationsIfNotCached(cache *pokecache.Cache, currentLocationID int, command string) {
	var wg sync.WaitGroup
	for i := 0; i < LocationCount; i++ {
		wg.Add(1)
		go cacheLocationIfNotCached(cache, currentLocationID, &wg)
		currentLocationID++
	}

	wg.Wait()
	return
}

func cacheLocationIfNotCached(cache *pokecache.Cache, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, ok := cache.Get(id); ok {
		return
	}

	location, pokemon := getPokeAPILocation(id)
	cache.Add(id, location, pokemon)
	return
}

func getPokeAPILocation(id int) (locationName string, pokemon []string) {
	address := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", id)
	response, errResponse := http.Get(address)
	if errResponse != nil {
		return fmt.Sprintf("Error retrieving location id %d: %w", id, errResponse), nil
	}
	defer response.Body.Close()

	body, errBody := ioutil.ReadAll(response.Body)
	if errBody != nil {
		return fmt.Sprintf("Unable to read body of location API response: %w", errBody),  nil
	}

	var locationResponse locationArea
	errUnmarshal := json.Unmarshal(body, &locationResponse)
	if errUnmarshal != nil {
		return "Location not available", nil  // TODO have some better output message here
	}

	return locationResponse.Name, getPokemonInLocation(locationResponse)
}

func getPokemonInLocation(location locationArea) (pokemonNames []string) {
	for _, encounter := range location.PokemonEncounters {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}

	return pokemonNames
}