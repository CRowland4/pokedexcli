package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"slices"
	"strings"
	"github.com/CRowland4/pokedexcli/internal/pokeapi"
	"github.com/CRowland4/pokedexcli/internal/pokecache"
)
const (
	lineSeparator = "\n\n+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+\n\n"
	helpMessage = "Usage:\nhelp: Display all commands\nmap: Display next 20 locations\nmapb: Display previous 20 locations\nexit: Exit the Pokedex"
	welcomMessage = "Welcome to the Pokedex!\n\nUsage:\nhelp: Display all commands\nexit: Exit the Pokedex"
)  // TODO add explore to helpMessage

func main() {
	fmt.Print(welcomMessage)

	cache := pokecache.NewCache(5 * time.Minute)
	locationCacher := pokeapi.LocationCacher()

	var currentLocations [pokeapi.LocationCount]string
	var currentPokemon []string

	for {
		fmt.Print(lineSeparator)
		command := getCommand()
		if command == "exit" {
			return
		} else if command == "help" {
			fmt.Print(helpMessage)
		} else if command == "map" || command == "mapb" {
			currentLocations = locationCacher(&cache, command)
			printLocations(currentLocations)
		} else if strings.Contains(command, "explore") {
			currentPokemon = cacheAreaPokemon(cache, command, currentLocations)
			printPokemon(currentPokemon)
		} else if strings.Contains(command, "catch") {
			catchPokemon(cache, command, currentPokemon)  // TODO create this function
		} else {
			fmt.Print("Command not recognized")
		}
	}

	return
}

func catchPokemon(cache pokecache.Cache, command string, currentPokemon []string) bool {
	return true  // TODO complete this function
}

func printPokemon(pokemon []string) {
	if len(pokemon) == 0 {
		return
	}

	fmt.Println("Found Pokemon:")
	for i := range pokemon {
		fmt.Println("	-", pokemon[i])
	}

	return
}

func cacheAreaPokemon(cache pokecache.Cache, command string, currentLocations [pokeapi.LocationCount]string) (pokemon []string) {
	commandPieces := strings.Split(command, " ")
	if len(commandPieces) != 2 {
		fmt.Print("Usage: explore <area-name>")
		return
	}

	if !slices.Contains(currentLocations[:], commandPieces[1]) {
		fmt.Print("You're not in this area right now!")
		return
	}

	fmt.Printf("\nExploring %s...\n", commandPieces[1])
	pokemon = getLocationPokemon(commandPieces[1], cache)
	// TODO cache pokemon
	return pokemon
}

func getLocationPokemon(location string, cache pokecache.Cache) (pokemon []string) {
	for _, entry := range cache.Info {
		if location == entry.LocationName {
			return entry.Pokemon
		}
	}

	return nil
}

func printLocations(locations [20]string) {
	for i := range locations {
		if locations[i] != "" { fmt.Println(locations[i]) }
	}

	return
}

func getCommand() (command string) {
	fmt.Print("Pokedex > ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
