package alecaframe

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	userHash := "test-hash"
	publicToken := "test-token"

	client := NewClient(userHash, publicToken)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.userHash != userHash {
		t.Errorf("Expected userHash to be %q, got %q", userHash, client.userHash)
	}

	if client.publicToken != publicToken {
		t.Errorf("Expected publicToken to be %q, got %q", publicToken, client.publicToken)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}
}

func TestRelicStructures(t *testing.T) {
	// Test that we can create and populate relic structures
	relic := Relic{
		Name:        "A5",
		Era:         "Axi",
		Tier:        "Common",
		Count:       10,
		IsVaulted:   true,
		Intact:      2,
		Radiant:     3,
		Flawless:    2,
		Exceptional: 3,
	}

	if relic.Name != "A5" {
		t.Errorf("Expected name to be A5, got %s", relic.Name)
	}

	if relic.Era != "Axi" {
		t.Errorf("Expected era to be Axi, got %s", relic.Era)
	}

	if !relic.IsVaulted {
		t.Error("Expected relic to be vaulted")
	}

	totalRefinements := relic.Intact + relic.Exceptional + relic.Flawless + relic.Radiant
	if totalRefinements != relic.Count {
		t.Errorf("Expected refinement counts to sum to %d, got %d", relic.Count, totalRefinements)
	}
}

func TestRelicInventoryStructure(t *testing.T) {
	inventory := RelicInventory{
		Relics: []Relic{
			{Name: "A5", Era: "Axi", Count: 10},
			{Name: "N3", Era: "Neo", Count: 5},
			{Name: "M2", Era: "Meso", Count: 3},
		},
	}

	if len(inventory.Relics) != 3 {
		t.Errorf("Expected 3 relics, got %d", len(inventory.Relics))
	}

	totalCount := 0
	for _, relic := range inventory.Relics {
		totalCount += relic.Count
	}

	if totalCount != 18 {
		t.Errorf("Expected total count to be 18, got %d", totalCount)
	}
}

func TestUserStatsStructure(t *testing.T) {
	stats := UserStats{
		TotalTrades:    100,
		PlatinumEarned: 5000,
		PlatinumSpent:  2000,
	}

	if stats.TotalTrades != 100 {
		t.Errorf("Expected 100 trades, got %d", stats.TotalTrades)
	}

	netPlatinum := stats.PlatinumEarned - stats.PlatinumSpent
	if netPlatinum != 3000 {
		t.Errorf("Expected net platinum to be 3000, got %d", netPlatinum)
	}
}

// Note: Actual API calls require valid credentials and are tested separately
// These tests focus on structure validation and client initialization
