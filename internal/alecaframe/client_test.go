package alecaframe

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	publicToken := "test-token"

	client := NewClient(publicToken)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
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
		Type:       "Axi",
		Refinement: "Radiant",
		Name:       "A5",
		Count:      10,
	}

	if relic.Name != "A5" {
		t.Errorf("Expected name to be A5, got %s", relic.Name)
	}

	if relic.Type != "Axi" {
		t.Errorf("Expected type to be Axi, got %s", relic.Type)
	}

	if relic.Refinement != "Radiant" {
		t.Errorf("Expected refinement to be Radiant, got %s", relic.Refinement)
	}

	if relic.Count != 10 {
		t.Errorf("Expected count to be 10, got %d", relic.Count)
	}
}

func TestRelicInventoryStructure(t *testing.T) {
	inventory := RelicInventory{
		Relics: []Relic{
			{Type: "Axi", Name: "A5", Refinement: "Intact", Count: 10},
			{Type: "Neo", Name: "N3", Refinement: "Radiant", Count: 5},
			{Type: "Meso", Name: "M2", Refinement: "Flawless", Count: 3},
		},
	}

	if len(inventory.Relics) != 3 {
		t.Errorf("Expected 3 relics, got %d", len(inventory.Relics))
	}

	totalCount := uint32(0)
	for _, relic := range inventory.Relics {
		totalCount += relic.Count
	}

	if totalCount != 18 {
		t.Errorf("Expected total count to be 18, got %d", totalCount)
	}
}

// Note: Actual API calls require valid credentials and are tested separately
// These tests focus on structure validation and client initialization
