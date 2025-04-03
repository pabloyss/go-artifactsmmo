package main

import (
	"fmt"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

// API-key
const APIKEY = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6InBhYmxveXVyaXNzQGdtYWlsLmNvbSIsInBhc3N3b3JkX2NoYW5nZWQiOiIifQ.NkQ7MduReGrZSyFxjOJ4dFfMkbGXrb-c8GdrkTM_cgQ"

// Characters
const fighterChar = "VespaMine"

func main() {

	client := artifactsmmo.NewClient(APIKEY, fighterChar)
	character, err := client.GetCharacterInfo(fighterChar)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("Welcome, %s (XP: %d/%d)!\nCurrent map: [%d,%d]\n", character.Name, character.Xp, character.MaxXp, character.X, character.Y)

	var response *models.DataSchemaCharacterAction

	// Move to gather zone
	_, err = client.Move(2, 0)
	if err != nil {
		fmt.Println("Error moving to gathering zone:", err)
	}

	// Game Loop
	for {

		// Gather
		response, err = client.Gather()
		if err != nil {
			fmt.Println("Error gathering:", err)
		} else {
			character = &response.Character
			fmt.Printf("Mining Lvl (XP: %d/%d):\n", character.MiningXp, character.MiningMaxXp)
		}

		inventoryFull := false

		// Check inventory
		for _, item := range character.Inventory {
			if item.Code == "copper_ore" && item.Quantity > 80 {
				inventoryFull = true
				break
			}
		}

		// Do routine for inventory full - crafting and depositing

		if inventoryFull {
			// Move to the Cooking
			_, err = client.Move(1, 5)
			if err != nil {
				fmt.Println("Error moving to the Mining Refinery:", err)
				break
			}

			// Craft potions
			response, err = client.Craft("copper", 8)
			if err != nil {
				fmt.Println("Error getting food:", err)
				break
			} else {
				character = &response.Character
				fmt.Printf("Mining Lvl (XP: %d/%d):\n", character.MiningXp, character.MiningMaxXp)
			}

			// Move to the bank
			_, err = client.Move(4, 1)
			if err != nil {
				fmt.Println("Error moving to the bank:", err)
				break
			}

			// Deposit all items
			for _, item := range character.Inventory {
				_, err = client.DepositBank(item.Code, item.Quantity)
				if err != nil {
					fmt.Println("Error depositing item %s:", item.Code, err)
					break
				}
			}

			// Move to gather zone
			_, err = client.Move(2, 0)
			if err != nil {
				fmt.Println("Error moving to gathering zone:", err)
				break
			}

		}
	}

}
