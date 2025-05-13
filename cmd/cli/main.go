package main

import (
	"fmt"
	"trading-bot/internal/interfaces/cli"
)

// main is the entrypoint for the CLI application.
func main() {

	// Defer a function to recover from a panic and handle errors gracefully.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("An error occurred:", r)
		}
	}()

	cli.ExecuteCLI()
}
