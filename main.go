package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/paperfeed/warframe-mcp/internal/alecaframe"
	"github.com/paperfeed/warframe-mcp/internal/config"
	"github.com/paperfeed/warframe-mcp/internal/market"
	"github.com/paperfeed/warframe-mcp/internal/recommender"
	"github.com/paperfeed/warframe-mcp/internal/worldstate"
)

const (
	serverName    = "warframe-mcp"
	serverVersion = "0.1.0"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Initialize API clients
	var alecaframeClient *alecaframe.Client
	if cfg.Alecaframe.PublicToken != "" {
		alecaframeClient = alecaframe.NewClient(cfg.Alecaframe.PublicToken)
	}

	worldStateClient := worldstate.NewClient(cfg.Cache.GetWorldStateTTL())
	marketClient := market.NewClient(cfg.Cache.GetMarketDataTTL())

	// Initialize recommender
	rec := recommender.NewRecommender(alecaframeClient, worldStateClient, marketClient)

	// Create MCP server
	s := server.NewMCPServer(serverName, serverVersion)

	// Register tools
	registerTools(s, rec, alecaframeClient, worldStateClient, marketClient)

	// Start server with stdio transport
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("serving: %w", err)
	}

	return nil
}

func loadConfig() (*config.Config, error) {
	// Start with default config
	cfg := &config.Config{
		Cache: config.CacheConfig{
			WorldStateTTL: 300,
			MarketDataTTL: 600,
		},
	}

	// Check for environment variables first (takes precedence)
	publicToken := os.Getenv("ALECAFRAME_PUBLIC_TOKEN")

	if publicToken != "" {
		cfg.Alecaframe.PublicToken = publicToken
		log.Println("Using Alecaframe credentials from environment variables")
		return cfg, nil
	}

	// Try to load config.json from current directory or home directory
	paths := []string{
		"config.json",
		os.ExpandEnv("$HOME/.config/warframe-mcp/config.json"),
		"/etc/warframe-mcp/config.json",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			fileCfg, err := config.Load(path)
			if err != nil {
				return nil, err
			}
			log.Printf("Loaded config from %s\n", path)
			return fileCfg, nil
		}
	}

	// No config file or env vars found
	log.Println("Warning: No config file or environment variables found, using defaults (Alecaframe features disabled)")
	return cfg, nil
}

func registerTools(s *server.MCPServer, rec *recommender.Recommender, aClient *alecaframe.Client, wsClient *worldstate.Client, mClient *market.Client) {
	// Tool: Get prime farming recommendations
	s.AddTool(mcp.Tool{
		Name:        "recommend_prime_farming",
		Description: "Get intelligent recommendations for which prime parts to farm based on current world state, vault status, and your relic inventory",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		recommendation, err := rec.RecommendPrimeParts()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		data, _ := json.MarshalIndent(recommendation, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// Tool: Get standing recommendations
	s.AddTool(mcp.Tool{
		Name:        "recommend_standing_farming",
		Description: "Get recommendations for which syndicate standing to farm based on active missions, world cycles, and current events",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		recommendation, err := rec.RecommendStanding()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		data, _ := json.MarshalIndent(recommendation, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// Tool: Get credit farming recommendations
	s.AddTool(mcp.Tool{
		Name:        "recommend_credit_farming",
		Description: "Get the best methods to farm credits based on current events and active bonuses",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		recommendation, err := rec.RecommendCredits()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		data, _ := json.MarshalIndent(recommendation, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// Tool: Get current world state
	s.AddTool(mcp.Tool{
		Name:        "get_world_state",
		Description: "Get the current Warframe world state including active alerts, invasions, fissures, sorties, and events",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		ws, err := wsClient.GetWorldState()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		// Create a summary instead of returning all data
		summary := buildWorldStateSummary(ws)
		return mcp.NewToolResultText(summary), nil
	})

	// Tool: Get Baro Ki'Teer status
	s.AddTool(mcp.Tool{
		Name:        "get_baro_status",
		Description: "Get information about Baro Ki'Teer (Void Trader) including his current location and inventory",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		ws, err := wsClient.GetWorldState()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		if ws.VoidTrader == nil {
			return mcp.NewToolResultText("No Void Trader information available"), nil
		}

		data, _ := json.MarshalIndent(ws.VoidTrader, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// Tool: Get Varzia status (Prime Resurgence)
	s.AddTool(mcp.Tool{
		Name:        "get_varzia_status",
		Description: "Get information about Varzia and Prime Resurgence, including currently available primes and rotation schedule",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		vaultTrader, err := wsClient.GetVaultTrader()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		data, _ := json.MarshalIndent(vaultTrader, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	// Tool: Get market price for an item
	s.AddTool(mcp.Tool{
		Name:        "get_market_price",
		Description: "Get the current market price and trading information for a Warframe item from warframe.market",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"item_name"},
			Properties: map[string]interface{}{
				"item_name": map[string]interface{}{
					"type":        "string",
					"description": "The name of the item to look up (e.g., 'rhino prime set', 'maiming strike')",
				},
			},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		itemName, ok := args["item_name"].(string)
		if !ok {
			return mcp.NewToolResultError("item_name must be a string"), nil
		}

		// Convert item name to URL format
		urlName := market.URLEncode(itemName)

		priceInfo, err := mClient.GetItemPrice(urlName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		summary := fmt.Sprintf(`Item: %s

Sell Orders: %d active
  Average: %.0f platinum
  Minimum: %d platinum

Buy Orders: %d active
  Average: %.0f platinum
  Maximum: %d platinum

Recommendation: %s`,
			itemName,
			priceInfo.TotalSellOrders,
			priceInfo.AvgSellPrice,
			priceInfo.MinSellPrice,
			priceInfo.TotalBuyOrders,
			priceInfo.AvgBuyPrice,
			priceInfo.MaxBuyPrice,
			getMarketRecommendation(priceInfo),
		)

		return mcp.NewToolResultText(summary), nil
	})

	// Tool: Search for items on the market
	s.AddTool(mcp.Tool{
		Name:        "search_market_items",
		Description: "Search for tradeable items on warframe.market by name",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"query"},
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query (e.g., 'prime', 'rhino', 'mod')",
				},
			},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		query, ok := args["query"].(string)
		if !ok {
			return mcp.NewToolResultError("query must be a string"), nil
		}

		items, err := mClient.SearchItems(query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		if len(items) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No items found for query: %s", query)), nil
		}

		// Limit to top 20 results
		if len(items) > 20 {
			items = items[:20]
		}

		result := fmt.Sprintf("Found %d items:\n\n", len(items))
		for i, item := range items {
			result += fmt.Sprintf("%d. %s\n", i+1, item.ItemName)
		}

		return mcp.NewToolResultText(result), nil
	})

	// Tool: Get user's relic inventory (Alecaframe)
	s.AddTool(mcp.Tool{
		Name:        "get_my_relics",
		Description: "Get your personal relic inventory from Alecaframe. Requires Alecaframe credentials to be configured. Shows all relics you own with counts and refinement levels.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
		if aClient == nil {
			return mcp.NewToolResultError("Alecaframe is not configured. Please set ALECAFRAME_PUBLIC_TOKEN environment variable or create a config.json file."), nil
		}

		inventory, err := aClient.GetRelicInventory()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error fetching relic inventory: %v", err)), nil
		}

		if len(inventory.Relics) == 0 {
			return mcp.NewToolResultText("No relics found in your inventory."), nil
		}

		// Build a nice summary
		result := fmt.Sprintf("=== Your Relic Inventory ===\n\nTotal relic entries: %d\n\n", len(inventory.Relics))

		// Group by type (era)
		byType := make(map[string][]alecaframe.Relic)
		for _, relic := range inventory.Relics {
			byType[relic.Type] = append(byType[relic.Type], relic)
		}

		// Display in order: Requiem, Axi, Neo, Meso, Lith
		types := []string{"Requiem", "Axi", "Neo", "Meso", "Lith"}
		for _, relicType := range types {
			relics, ok := byType[relicType]
			if !ok || len(relics) == 0 {
				continue
			}

			result += fmt.Sprintf("## %s Relics (%d)\n", relicType, len(relics))
			for _, relic := range relics {
				result += fmt.Sprintf("  %s %s (%s): %d\n", relic.Type, relic.Name, relic.Refinement, relic.Count)
			}
			result += "\n"
		}

		return mcp.NewToolResultText(result), nil
	})
}

func buildWorldStateSummary(ws *worldstate.WorldState) string {
	summary := "=== Current World State ===\n\n"

	// Alerts
	activeAlerts := 0
	for _, alert := range ws.Alerts {
		if alert.Active {
			activeAlerts++
		}
	}
	if activeAlerts > 0 {
		summary += fmt.Sprintf("🚨 Alerts: %d active\n", activeAlerts)
	}

	// Arbitration
	if ws.Arbitration != nil && ws.Arbitration.Active {
		summary += fmt.Sprintf("⚔️  Arbitration: %s (%s) on %s\n", ws.Arbitration.Type, ws.Arbitration.Enemy, ws.Arbitration.Node)
	}

	// Sortie
	if ws.Sortie != nil && ws.Sortie.Active {
		summary += fmt.Sprintf("🎯 Sortie: %s (%s) - %d missions\n", ws.Sortie.Boss, ws.Sortie.Faction, len(ws.Sortie.Variants))
	}

	// Fissures
	activeFissures := 0
	for _, fissure := range ws.Fissures {
		if fissure.Active {
			activeFissures++
		}
	}
	if activeFissures > 0 {
		summary += fmt.Sprintf("🌀 Fissures: %d active\n", activeFissures)
	}

	// Invasions
	activeInvasions := 0
	for _, invasion := range ws.Invasions {
		if !invasion.Completed {
			activeInvasions++
		}
	}
	if activeInvasions > 0 {
		summary += fmt.Sprintf("⚔️  Invasions: %d active\n", activeInvasions)
	}

	// Baro
	if ws.VoidTrader != nil {
		if ws.VoidTrader.Active {
			summary += fmt.Sprintf("💎 Baro Ki'Teer: At %s with %d items\n", ws.VoidTrader.Location, len(ws.VoidTrader.Inventory))
		} else {
			summary += "💎 Baro Ki'Teer: Not currently visiting\n"
		}
	}

	// Cycles
	summary += "\n=== Open World Cycles ===\n\n"

	if ws.CetusCycle != nil {
		state := "Day ☀️"
		if !ws.CetusCycle.IsDay {
			state = "Night 🌙"
		}
		summary += fmt.Sprintf("Plains of Eidolon: %s (%s remaining)\n", state, ws.CetusCycle.TimeLeft)
	}

	if ws.VallisEarth != nil {
		state := "Warm"
		if !ws.VallisEarth.IsWarm {
			state = "Cold"
		}
		summary += fmt.Sprintf("Orb Vallis: %s (%s remaining)\n", state, ws.VallisEarth.TimeLeft)
	}

	if ws.CambionDrift != nil {
		summary += fmt.Sprintf("Cambion Drift: %s\n", ws.CambionDrift.State)
	}

	return summary
}

func getMarketRecommendation(info *market.PriceInfo) string {
	if info.TotalSellOrders == 0 && info.TotalBuyOrders == 0 {
		return "No active orders - item may not be tradeable or is very rare"
	}

	if info.TotalSellOrders < 5 {
		return "Low supply - good item to farm and sell"
	}

	if info.AvgSellPrice > 100 {
		return "High value item - worth farming for platinum"
	}

	if info.AvgSellPrice < 10 {
		return "Low value item - not recommended for platinum farming"
	}

	return "Moderate value - check if you have spare sets to sell"
}
