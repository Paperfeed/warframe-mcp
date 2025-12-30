package main

import (
	"fmt"
	"os"

	"github.com/paperfeed/warframe-mcp/internal/alecaframe"
)

func main() {
	publicToken := os.Getenv("ALECAFRAME_PUBLIC_TOKEN")

	if publicToken == "" {
		fmt.Println("Error: ALECAFRAME_PUBLIC_TOKEN environment variable must be set")
		os.Exit(1)
	}

	fmt.Printf("Testing Alecaframe API with credentials...\n")
	fmt.Printf("Public Token: %s...\n\n", maskString(publicToken, 8))

	client := alecaframe.NewClient(publicToken)

	// Test 1: Get Relic Inventory
	fmt.Println("=== Test 1: Get Relic Inventory ===")
	inventory, err := client.GetRelicInventory()
	if err != nil {
		fmt.Printf("❌ Error: %v\n\n", err)
	} else {
		fmt.Printf("✅ Success! Found %d relic entries\n", len(inventory.Relics))
		if len(inventory.Relics) > 0 {
			fmt.Println("\nFirst 5 relics:")
			for i, relic := range inventory.Relics {
				if i >= 5 {
					break
				}
				fmt.Printf("  %s %s (%s): %d\n", relic.Type, relic.Name, relic.Refinement, relic.Count)
			}
		}
		fmt.Println()
	}

	// Test 2: Get User Stats
	fmt.Println("=== Test 2: Get User Stats ===")
	stats, err := client.GetUserStats()
	if err != nil {
		fmt.Printf("❌ Error: %v\n\n", err)
	} else {
		fmt.Printf("✅ Success!\n")
		fmt.Printf("Total Trades: %d\n", stats.TotalTrades)
		fmt.Printf("Platinum Earned: %d\n", stats.PlatinumEarned)
		fmt.Printf("Platinum Spent: %d\n", stats.PlatinumSpent)
		fmt.Printf("Net Platinum: %d\n", stats.PlatinumEarned-stats.PlatinumSpent)
		fmt.Println()
	}
}

func maskString(s string, showChars int) string {
	if len(s) <= showChars {
		return s
	}
	return s[:showChars] + "..."
}
