package main

import (
	"fmt"
	"bufio"
	"os"
	"github.com/CRowland4/pokedexcli/internal/pokeapi"
)

func main() {
	fmt.Println("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\nhelp: Display all commands\nexit: Exit the Pokedex")

	locationGetter := pokeapi.LocationGetter()
	for {
		fmt.Print("\n\n+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+\n\n")
		switch command := getCommand(); command {
		case "exit":
			return
		case "help":
			fmt.Print("Usage:\nhelp: Display all commands\nmap: Display next 20 locations\nmapb: Display previous 20 locations\nexit: Exit the Pokedex")  // TODO consolodate into variable
		case "map", "mapb":
			locations, err := locationGetter(command)
			printLocations(locations, err)
		default:
			fmt.Print("Command not recognized")
		}
	}
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