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
	helpMessage = `Usage:
	help: Display all commands
	map: Display next 20 locations
	mapb: Display previous 20 locations
	explore <area>: Discover the pokemon located in one of your current locations
	catch <pokemon>: Attempt to catch one of the pokemon you have discovered form exploring an area
	inspect <pokemon>: Inspect a pokemon that you have caught
	pokedex: View the names of all the pokemon that you have caught
	exit: Exit the Pokedex`
	welcomMessage = "Welcome to the Pokedex!\n\nUsage:\nhelp: Display all commands\nexit: Exit the Pokedex"
)

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
		} else if strings.Contains(command, "inspect") {
			inspectPokemon(cache, command)
		} else if command == "pokedex" {
			printCaughtPokemon(cache)
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
	count := 0
	for _, name := range locations {
		if name != "" { count++ }
	}

	if count == 0 {
		fmt.Println("Nothing to explore here...")
		return
	}

	for _, location := range locations {
		if location != "" { fmt.Println(location) }
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
	baseExperience := cache.Pokemon[pokemonToCatch].BaseExperience

	if rand.Intn(100000) > baseExperience {
		fmt.Println(pokemonToCatch, "was caught!")
		cachedPokemon := cache.Pokemon[pokemonToCatch]
		cachedPokemon.IsCaught = true
		cache.Pokemon[pokemonToCatch] = cachedPokemon  // Make a cache method for this, since it modifies the cache?
	} else {
		fmt.Println(pokemonToCatch, "escaped!")
	}

	return
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
// inspect command

func inspectPokemon(cache pokecache.Cache, command string) {
	commandPieces := strings.Split(command, " ")
	if len(commandPieces) != 2 {
		fmt.Println("Usage: inspect <name of pokemon>")
		return
	}

	pokemonToInspect := commandPieces[1]
	if _, ok := cache.Pokemon[pokemonToInspect]; !ok {
		fmt.Println("You haven't discovered a pokemon named", pokemonToInspect, "yet!")
		return
	}

	if !cache.Pokemon[pokemonToInspect].IsCaught {
		fmt.Println("You haven't caught a", pokemonToInspect, "yet!")
		return
	}

	printPokemonInformation(cache, pokemonToInspect)
	return
}

func printPokemonInformation(cache pokecache.Cache, pokemon string) {
	fmt.Println("Name:", pokemon)
	fmt.Println("Height:", cache.Pokemon[pokemon].Height)
	fmt.Println("Weight:", cache.Pokemon[pokemon].Weight)
	fmt.Println("Stats:")
	fmt.Println("  -hp:", cache.Pokemon[pokemon].HP)
	fmt.Println("  -attack:", cache.Pokemon[pokemon].Attack)
	fmt.Println("  -defense:", cache.Pokemon[pokemon].Defense)
	fmt.Println("  -special-attack:", cache.Pokemon[pokemon].SpecialAttack)
	fmt.Println("  -special-defense:", cache.Pokemon[pokemon].SpecialDefense)
	fmt.Println("  -Speed:", cache.Pokemon[pokemon].Speed)
	fmt.Println("Types:")

	for _, type_ := range cache.Pokemon[pokemon].Types {
		fmt.Println("  -" + type_)
	}

	return
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
// pokedex command

func printCaughtPokemon(cache pokecache.Cache) {
	caughtPokemon := cache.GetCaughtPokemon()

	if len(caughtPokemon) == 0 {
		fmt.Println("You haven't caught any pokemon yet!")
		return
	}

	fmt.Println("Your Pokedex:")
	for _, name := range caughtPokemon {
		fmt.Println("  -", name)
	}

	return
}