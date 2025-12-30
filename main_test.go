package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paperfeed/warframe-mcp/internal/market"
	"github.com/paperfeed/warframe-mcp/internal/worldstate"
)

func TestLoadConfigWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("ALECAFRAME_PUBLIC_TOKEN", "test-token-env")
	defer os.Unsetenv("ALECAFRAME_PUBLIC_TOKEN")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	if cfg.Alecaframe.PublicToken != "test-token-env" {
		t.Errorf("Expected publicToken from env to be 'test-token-env', got %q", cfg.Alecaframe.PublicToken)
	}
}

func TestLoadConfigWithFile(t *testing.T) {
	// Ensure no env vars are set
	os.Unsetenv("ALECAFRAME_PUBLIC_TOKEN")

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"alecaframe": {
			"publicToken": "test-token-file"
		},
		"cache": {
			"worldStateTTL": 300,
			"marketDataTTL": 600
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Change to temp directory so config.json is found
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	if cfg.Alecaframe.PublicToken != "test-token-file" {
		t.Errorf("Expected publicToken from file to be 'test-token-file', got %q", cfg.Alecaframe.PublicToken)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Ensure no env vars or config files
	os.Unsetenv("ALECAFRAME_PUBLIC_TOKEN")

	// Change to a temp directory with no config
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	// Defaults should be set
	if cfg.Cache.WorldStateTTL != 300 {
		t.Errorf("Expected default WorldStateTTL to be 300, got %d", cfg.Cache.WorldStateTTL)
	}

	if cfg.Cache.MarketDataTTL != 600 {
		t.Errorf("Expected default MarketDataTTL to be 600, got %d", cfg.Cache.MarketDataTTL)
	}

	// Alecaframe should be empty
	if cfg.Alecaframe.PublicToken != "" {
		t.Errorf("Expected empty publicToken with defaults, got %q", cfg.Alecaframe.PublicToken)
	}
}

func TestLoadConfigEnvPrecedence(t *testing.T) {
	// Create a config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"alecaframe": {
			"publicToken": "test-token-file"
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set environment variables (should take precedence)
	os.Setenv("ALECAFRAME_PUBLIC_TOKEN", "test-token-env")
	defer os.Unsetenv("ALECAFRAME_PUBLIC_TOKEN")

	// Change to temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	// Environment variables should win
	if cfg.Alecaframe.PublicToken != "test-token-env" {
		t.Errorf("Expected env vars to take precedence, got publicToken %q", cfg.Alecaframe.PublicToken)
	}
}

func TestBuildWorldStateSummary(t *testing.T) {
	ws := &worldstate.WorldState{
		Timestamp: "2024-01-01T00:00:00.000Z",
		Fissures: []worldstate.Fissure{
			{Active: true},
			{Active: true},
			{Active: false},
		},
		Sortie: &worldstate.Sortie{
			Active:   true,
			Boss:     "The Sergeant",
			Faction:  "Corpus",
			Variants: []worldstate.SortieVariant{{}, {}},
		},
		Arbitration: &worldstate.Arbitration{
			Active: true,
			Node:   "Mot (Void)",
			Type:   "Survival",
			Enemy:  "Corrupted",
		},
		CetusCycle: &worldstate.CetusCycle{
			IsDay:    true,
			State:    "day",
			TimeLeft: "2h 15m",
		},
	}

	summary := buildWorldStateSummary(ws)

	if summary == "" {
		t.Error("Expected summary to be generated, got empty string")
	}

	t.Logf("Generated world state summary:\n%s", summary)

	// Check for expected content
	if !contains(summary, "Fissures") {
		t.Error("Expected summary to mention fissures")
	}

	if !contains(summary, "Sortie") {
		t.Error("Expected summary to mention sortie")
	}

	if !contains(summary, "Arbitration") {
		t.Error("Expected summary to mention arbitration")
	}
}

func TestGetMarketRecommendation(t *testing.T) {
	tests := []struct {
		name     string
		info     *market.PriceInfo
		expected string
	}{
		{
			name: "No orders",
			info: &market.PriceInfo{
				TotalSellOrders: 0,
				TotalBuyOrders:  0,
			},
			expected: "No active orders",
		},
		{
			name: "Low supply",
			info: &market.PriceInfo{
				TotalSellOrders: 3,
				AvgSellPrice:    50,
			},
			expected: "Low supply",
		},
		{
			name: "High value",
			info: &market.PriceInfo{
				TotalSellOrders: 10,
				AvgSellPrice:    150,
			},
			expected: "High value",
		},
		{
			name: "Low value",
			info: &market.PriceInfo{
				TotalSellOrders: 20,
				AvgSellPrice:    5,
			},
			expected: "Low value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMarketRecommendation(tt.info)
			if !contains(result, tt.expected) {
				t.Errorf("Expected recommendation to contain %q, got %q", tt.expected, result)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
