package main

import (
	"fmt"

	"github.com/0xN0x/go-artifactsmmo"
	"github.com/0xN0x/go-artifactsmmo/models"
)

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
	}

	fmt.Printf("Welcome, %s (XP: %d/%d)!\nCurrent map: [%d,%d]\n", character.Name, character.Xp, character.MaxXp, character.X, character.Y)

	var response *models.DataSchemaCharacterAction

	// Move to fight zone (red slime)
	_, err = client.Move(4, -1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Game Loop
	for {

		switch {
		case character.Hp < character.MaxHp-75:
			// Check presence of item
			hasFood := false

			// Run through the inventory and check if there is food
			for _, item := range character.Inventory {
				if item.Code == "cooked_gudgeon" && item.Quantity > 0 {
					hasFood = true
					break
				}
			}

			if !hasFood {
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

				// Withdraw Potions
				response, err = client.WithdrawBank("small_health_potion", 40)
				if err != nil {
					fmt.Println("Error getting potions:", err)
					break
				} else {
					character = &response.Character
				}

				// Equip potions
				response, err = client.Equip("small_health_potion", models.Utility1, 40)
				if err != nil {
					fmt.Println("Error equiping potion:", err)
					break
				} else {
					character = &response.Character
				}

				// Get Food
				response, err = client.WithdrawBank("cooked_gudgeon", 50)
				if err != nil {
					fmt.Println("Error getting food:", err)
					break
				} else {
					character = &response.Character
				}

				// Move back to slime
				_, err = client.Move(4, -1)
				if err != nil {
					fmt.Println("Error moving after the bank:", err)
					break
				}

			}

			// Use Food
			response, err = client.UseItem("cooked_gudgeon", 1)
			if err != nil {
				// No Food
				fmt.Println("Error using food:", err)
			} else {
				character = &response.Character
			}

		case response != nil && (response.Fight.Result == "loss"):
			// Move back
			_, err = client.Move(4, -1)
			if err != nil {
				fmt.Println("Error moving after loss:", err)
			}

			// Rest
			_, err = client.Rest()
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				character = &response.Character
			}

		default:
			response, err = client.Fight()
			character = &response.Character
			if err != nil {
				fmt.Println("Error:", err)
			}

			fmt.Printf("%s (XP: %d/%d)!\nHP: [%d/%d]\n", response.Character.Name, response.Character.Xp, response.Character.MaxXp, response.Character.Hp, response.Character.MaxHp)
		}
	}

}
