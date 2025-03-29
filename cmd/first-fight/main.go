package main

import (
	"fmt"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
)

func usage() {
	fmt.Println("Usage: go run main.go <api-token> <character-name>")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	client := artifactsmmo.NewClient(APIKEY, fighterChar)
	character, err := client.GetCharacterInfo(fighterChar)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Welcome, %s (XP: %d/%d)!\nCurrent map: [%d,%d]\n", character.Name, character.Xp, character.MaxXp, character.X, character.Y)

	fight, err := client.Fight()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("[%s] Fight result: +%d xp, +%d gold\n", fight.Fight.Result, fight.Fight.Xp, fight.Fight.Gold)
}
