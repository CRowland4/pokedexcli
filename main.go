package main

import (
	"fmt"
	"bufio"
	"os"
)

var commands = map[string]string{
	"help": "Usage:\n\nhelp: Display all commands\nexit: Exit the Pokedex\n\n",
	"default": "Command not recognized\n",
}

func main() {
	fmt.Println("Welcome to the Pokedex!\n")
	fmt.Println("Usage:\nhelp: Display all commands\nexit: Exit the Pokedex\n\n")

	for {
		fmt.Println("+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+")
		command := getCommand()
		if command == "exit" {
			return
		}

		if response, ok := commands[command]; ok {
			fmt.Println(response)
		} else {
			fmt.Println(commands["default"])
		}
	}
	}

func getCommand() (command string) {
	fmt.Print("Pokedex > ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}