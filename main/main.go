package main

import (
	"fmt"
	"log"

	"github.com/KaylaHood/CigarBidFreefallBot/cigarbot"
)

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	// A slice of names
	names := []string{"Gladys", "Samantha", "Darrin"}

	// Get a greeting message and print it.
	messages, err := cigarbot.Hellos(names)
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		log.Fatal(err)
	}

	// If no error was caught, print message
	fmt.Println(messages)
}
