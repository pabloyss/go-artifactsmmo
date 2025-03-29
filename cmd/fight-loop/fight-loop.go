package main

import (
	"fmt"
	"os"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

func usage() {
	fmt.Println("Usage: go run main.go <api-token> <character-name>")
}

// Make the program sleep for n seconds

// API-key
const APIKEY = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6InBhYmxveXVyaXNzQGdtYWlsLmNvbSIsInBhc3N3b3JkX2NoYW5nZWQiOiIifQ.NkQ7MduReGrZSyFxjOJ4dFfMkbGXrb-c8GdrkTM_cgQ"

// Characters
const fighterChar = "Vespa"

func main() {

	client := artifactsmmo.NewClient(APIKEY, fighterChar)
	character, err := client.GetCharacterInfo(fighterChar)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Welcome, %s (XP: %d/%d)!\nCurrent map: [%d,%d]\n", character.Name, character.Xp, character.MaxXp, character.X, character.Y)

	var response *models.CharacterFight

	// Move to fight zone (red slime)
	_, err = client.Move(1, -1)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Game Loop
	for {

		switch {
		case response.Character.Hp < response.Character.MaxHp-75:
			// Use Food

			// No Food
			// Move to the bank
			// Get Itens

		case response.Fight.Result == "loss":
			// Move back

			// Rest

		default:
			// Move to fight zone (red slime)
			_, err = client.Move(1, -1)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			response, err = client.Fight()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			fmt.Printf("%s (XP: %d/%d)!\nHP: [%d/%d]\n", response.Character.Name, response.Character.Xp, response.Character.MaxXp, response.Character.Hp, response.Character.MaxHp)
		}

	}

}
