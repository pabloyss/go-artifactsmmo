package main

import (
	"fmt"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
)

func usage() {
	fmt.Println("Usage: go run main.go <api-token> <character-name>")
}

// API-key
const APIKEY = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6InBhYmxveXVyaXNzQGdtYWlsLmNvbSIsInBhc3N3b3JkX2NoYW5nZWQiOiIifQ.NkQ7MduReGrZSyFxjOJ4dFfMkbGXrb-c8GdrkTM_cgQ"

// Characters
const fighterChar = "Vespa"

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	client := artifactsmmo.NewClient(APIKEY, fighterChar)
	character, err := client.GetCharacterInfo(fighterChar)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Welcome, %s (XP: %d/%d)!\n", character.Name, character.Xp, character.MaxXp)
}
