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
	for {
		fmt.Print(lineSeparator)
		command := getCommand()
		if command == "exit" {
			return
		} else if command == "help" {
			fmt.Print(helpMessage)
		} else if command == "map" || command == "mapb" {
			currentLocations, err := locationCacher(&cache, command)
			printLocations(currentLocations, err)
		} else if strings.Contains(command, "explore") {
			exploreArea(cache, command, currentLocations)
		} else {
			fmt.Print("Command not recognized")
		}
	}

	return
}

func exploreArea(cache pokecache.Cache, command string, currentLocations [pokeapi.LocationCount]string) {
	commandPieces := strings.Split(command, " ")
	if len(commandPieces) != 2 {
		fmt.Print("Usage: explore <area-name>")
		return
	}

	if !slices.Contains(currentLocations[:], commandPieces[1]) {
		fmt.Print("You're not in this area right now!")
		return
	}

	fmt.Printf("Exploring %s...\n", commandPieces[1])
	fmt.Printf("Found Pokemon:")
	for _, pokemon := range getLocationPokemon(commandPieces[1], cache) {
		fmt.Printf("\n%s", pokemon)
	}

	return
}

func getLocationPokemon(location string, cache pokecache.Cache) (pokemon []string) {
	for _, entry := range cache.Info {
		if location == entry.LocationName {
			return entry.Pokemon
		}
	}

	return nil
}

func printLocations(locations [20]string, err error) {
	if err != nil {
		fmt.Printf("Error retrieving locations: %w\n\n", err)  // TODO print just message of error
		return
	}

	for i := range locations {
		fmt.Println(locations[i])
	}

	return
}

func getCommand() (command string) {
	fmt.Print("Pokedex > ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
