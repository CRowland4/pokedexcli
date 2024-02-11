package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"math/rand"
	"slices"
	"strings"
	"github.com/CRowland4/pokedexcli/internal/pokeapi"
	"github.com/CRowland4/pokedexcli/internal/pokecache"
)
const (
	lineSeparator = "\n\n+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+\n\n"
	helpMessage = "Usage:\nhelp: Display all commands\nmap: Display next 20 locations\nmapb: Display previous 20 locations\nexit: Exit the Pokedex"
	welcomMessage = "Welcome to the Pokedex!\n\nUsage:\nhelp: Display all commands\nexit: Exit the Pokedex"
)  // TODO add explore + catch to helpMessage

func main() {
	fmt.Print(welcomMessage)

	cache := pokecache.NewCache(5 * time.Minute)
	locationCacher := pokeapi.LocationCacher()

	var currentLocations [pokeapi.LocationCount]string
	var currentPokemon []string

	for {
		command := getCommand()
		if command == "exit" {
			return
		} else if command == "help" {
			fmt.Print(helpMessage)
		} else if command == "map" || command == "mapb" {
			currentLocations = locationCacher(&cache, command)
			currentPokemon = []string{}
			printLocations(currentLocations)
		} else if strings.Contains(command, "explore") {
			currentPokemon = getAreaPokemon(command, currentLocations, cache)
			printAreaPokemon(currentPokemon)
		} else if strings.Contains(command, "catch") {
			catchPokemon(cache, command, currentPokemon)
		} else {
			fmt.Print("Command not recognized")
		}
	}

	return
}

func getCommand() (command string) {
	fmt.Print(lineSeparator)
	fmt.Print("Pokedex > ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
// map & mapb commands
func printLocations(locations [pokeapi.LocationCount]string) {
	for i := range locations {
		if locations[i] != "" { fmt.Println(locations[i]) }
	}

	return
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
// explore command

func getAreaPokemon(command string, currentLocations [pokeapi.LocationCount]string, cache pokecache.Cache) (names []string) {
	commandPieces := strings.Split(command, " ")
	if len(commandPieces) != 2 {
		fmt.Print("Usage: explore <area-name>")
		return
	}

	location := commandPieces[1]
	if !slices.Contains(currentLocations[:], location) {
		fmt.Print("You're not in this area right now!")
		return
	}

	return getLocationPokemon(location, cache)
}

func getLocationPokemon(location string, cache pokecache.Cache) (pokemon []string) {
	for _, entry := range cache.Info {
		if location == entry.LocationName {
			return entry.LocationPokemon
		}
	}

	return nil
}

func printAreaPokemon(pokemon []string) {
	for _, name := range pokemon {
		fmt.Println("  -", name)
	}

	return
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
// catch command

func catchPokemon(cache pokecache.Cache, command string, currentPokemon []string) {
	commandPieces := strings.Split(command, " ")
	if len(commandPieces) != 2 {
		fmt.Println("Usage: catch <name of pokemon>")
		return
	}

	pokemonToCatch := commandPieces[1]
	if !slices.Contains(currentPokemon, pokemonToCatch) {
		fmt.Printf("%s isn't here!\n", pokemonToCatch)
		return
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonToCatch)
	baseExperience, _ := cache.GetPokemonBaseExperience(pokemonToCatch)

	if float64(rand.Intn(100000)) > baseExperience {  // TODO 100,000 is arbitrary and should be replaced
		fmt.Println(pokemonToCatch, "was caught!")
	} else {
		fmt.Println(pokemonToCatch, "escaped!")
	}

	return
}
