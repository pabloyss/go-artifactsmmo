package main

import (
	"fmt"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

// API-key
const APIKEY = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6InBhYmxveXVyaXNzQGdtYWlsLmNvbSIsInBhc3N3b3JkX2NoYW5nZWQiOiIifQ.NkQ7MduReGrZSyFxjOJ4dFfMkbGXrb-c8GdrkTM_cgQ"

// Characters
const fighterChar = "VespaAlchemy"

func main() {

	client := artifactsmmo.NewClient(APIKEY, fighterChar)
	character, err := client.GetCharacterInfo(fighterChar)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("Welcome, %s (XP: %d/%d)!\nCurrent map: [%d,%d]\n", character.Name, character.Xp, character.MaxXp, character.X, character.Y)

	var response *models.DataSchemaCharacterAction

	// Move to gather zone
	_, err = client.Move(2, 2)
	if err != nil {
		fmt.Println("Error moving to gathering zone:", err)
	}

	// Game Loop
	for {

		// Gather
		response, err = client.Gather()
		if err != nil {
			fmt.Println("Error gathering:", err)
			break
		} else {
			character = &response.Character
			fmt.Printf("Alchemy Lvl (XP: %d/%d):\n", character.AlchemyXp, character.AlchemyMaxXp)
		}

		inventoryFull := false

		// Check inventory
		for _, item := range character.Inventory {
			if item.Code == "sunflower" && item.Quantity > 99 {
				inventoryFull = true
				break
			}
		}

		// Do routine for inventory full - crafting and depositing

		if inventoryFull {
			// Move to the AlchemyCraft
			_, err = client.Move(2, 3)
			if err != nil {
				fmt.Println("Error moving to the AlchemyCraft:", err)
				break
			}

			// Craft potions
			response, err = client.Craft("small_health_potion", 33)
			if err != nil {
				fmt.Println("Error getting food:", err)
				break
			} else {
				character = &response.Character
				fmt.Printf("Alchemy Lvl (XP: %d/%d):", character.AlchemyXp, character.AlchemyMaxXp)
			}

			// Move to the bank
			_, err = client.Move(4, 1)
			if err != nil {
				fmt.Println("Error moving to the bank:", err)
				break
			}

			// Deposit potions
			response, err = client.DepositBank("small_health_potion", 33)
			if err != nil {
				fmt.Println("Error depositing potions:", err)
				break
			} else {
				character = &response.Character
			}

			// Move to gather zone
			_, err = client.Move(2, 2)
			if err != nil {
				fmt.Println("Error moving to gathering zone:", err)
				break
			}

		}
	}

}
