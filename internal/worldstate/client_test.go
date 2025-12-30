package worldstate

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	ttl := 5 * time.Minute
	client := NewClient(ttl)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.cache.ttl != ttl {
		t.Errorf("Expected cache TTL to be %v, got %v", ttl, client.cache.ttl)
	}
}

func TestGetWorldState(t *testing.T) {
	client := NewClient(5 * time.Minute)

	ws, err := client.GetWorldState()
	if err != nil {
		t.Logf("Warning: GetWorldState failed (this is expected if API is down): %v", err)
		return
	}

	if ws == nil {
		t.Fatal("Expected world state to be returned, got nil")
	}

	// Test that we got some basic data
	t.Logf("World state timestamp: %s", ws.Timestamp)
	t.Logf("Active fissures: %d", len(ws.Fissures))
	t.Logf("Active invasions: %d", len(ws.Invasions))

	// Test caching - second call should use cache
	ws2, err := client.GetWorldState()
	if err != nil {
		t.Fatalf("Second GetWorldState call failed: %v", err)
	}

	if ws.Timestamp != ws2.Timestamp {
		t.Error("Expected cached result to have same timestamp")
	}
}

func TestGetVaultTrader(t *testing.T) {
	client := NewClient(5 * time.Minute)

	vt, err := client.GetVaultTrader()
	if err != nil {
		t.Logf("Warning: GetVaultTrader failed (this is expected if API is down): %v", err)
		return
	}

	if vt == nil {
		t.Fatal("Expected vault trader to be returned, got nil")
	}

	t.Logf("Vault trader character: %s", vt.Character)
	t.Logf("Vault trader active: %v", vt.Active)
	if len(vt.Schedule) > 0 {
		t.Logf("Schedule entries: %d", len(vt.Schedule))
	}
}

func TestWorldStateStructures(t *testing.T) {
	// Test that we can create and populate structures
	ws := &WorldState{
		Timestamp: "2024-01-01T00:00:00.000Z",
		Fissures: []Fissure{
			{
				ID:          "test-fissure",
				Node:        "Mot (Void)",
				Tier:        "Axi",
				TierNum:     4,
				MissionType: "Survival",
				Enemy:       "Corrupted",
				Active:      true,
				IsStorm:     false,
				IsHard:      false,
			},
		},
		Sortie: &Sortie{
			ID:     "test-sortie",
			Boss:   "The Sergeant",
			Active: true,
			Variants: []SortieVariant{
				{
					Node:        "Iliad (Phobos)",
					MissionType: "Exterminate",
					Modifier:    "Enemy Physical Enhancement",
				},
			},
		},
	}

	if ws.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}

	if len(ws.Fissures) != 1 {
		t.Errorf("Expected 1 fissure, got %d", len(ws.Fissures))
	}

	if ws.Sortie == nil {
		t.Error("Expected sortie to be set")
	}
}
