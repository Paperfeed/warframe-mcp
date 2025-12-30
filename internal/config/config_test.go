package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"alecaframe": {
			"userHash": "test-hash",
			"publicToken": "test-token"
		},
		"cache": {
			"worldStateTTL": 600,
			"marketDataTTL": 1200
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify Alecaframe config
	if cfg.Alecaframe.UserHash != "test-hash" {
		t.Errorf("Expected userHash to be 'test-hash', got %q", cfg.Alecaframe.UserHash)
	}

	if cfg.Alecaframe.PublicToken != "test-token" {
		t.Errorf("Expected publicToken to be 'test-token', got %q", cfg.Alecaframe.PublicToken)
	}

	// Verify cache config
	if cfg.Cache.WorldStateTTL != 600 {
		t.Errorf("Expected worldStateTTL to be 600, got %d", cfg.Cache.WorldStateTTL)
	}

	if cfg.Cache.MarketDataTTL != 1200 {
		t.Errorf("Expected marketDataTTL to be 1200, got %d", cfg.Cache.MarketDataTTL)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Create a minimal config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"alecaframe": {
			"userHash": "test-hash",
			"publicToken": "test-token"
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify defaults are set
	if cfg.Cache.WorldStateTTL != 300 {
		t.Errorf("Expected default worldStateTTL to be 300, got %d", cfg.Cache.WorldStateTTL)
	}

	if cfg.Cache.MarketDataTTL != 600 {
		t.Errorf("Expected default marketDataTTL to be 600, got %d", cfg.Cache.MarketDataTTL)
	}
}

func TestLoadConfigInvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/config.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent config file")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Write invalid JSON
	err := os.WriteFile(configPath, []byte("{ invalid json }"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestGetWorldStateTTL(t *testing.T) {
	cfg := CacheConfig{
		WorldStateTTL: 300,
	}

	ttl := cfg.GetWorldStateTTL()
	expected := 300 * time.Second

	if ttl != expected {
		t.Errorf("Expected TTL to be %v, got %v", expected, ttl)
	}
}

func TestGetMarketDataTTL(t *testing.T) {
	cfg := CacheConfig{
		MarketDataTTL: 600,
	}

	ttl := cfg.GetMarketDataTTL()
	expected := 600 * time.Second

	if ttl != expected {
		t.Errorf("Expected TTL to be %v, got %v", expected, ttl)
	}
}
