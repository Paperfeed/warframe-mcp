package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds all configuration for the Warframe MCP server
type Config struct {
	Alecaframe AlecaframeConfig `json:"alecaframe"`
	Cache      CacheConfig      `json:"cache"`
}

// AlecaframeConfig holds Alecaframe API credentials
type AlecaframeConfig struct {
	PublicToken string `json:"publicToken"`
}

// CacheConfig holds cache TTL settings
type CacheConfig struct {
	WorldStateTTL  int `json:"worldStateTTL"`  // seconds
	MarketDataTTL  int `json:"marketDataTTL"`  // seconds
}

// Load reads configuration from a JSON file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set defaults if not specified
	if cfg.Cache.WorldStateTTL == 0 {
		cfg.Cache.WorldStateTTL = 300 // 5 minutes
	}
	if cfg.Cache.MarketDataTTL == 0 {
		cfg.Cache.MarketDataTTL = 600 // 10 minutes
	}

	return &cfg, nil
}

// GetWorldStateTTL returns the world state cache duration
func (c *CacheConfig) GetWorldStateTTL() time.Duration {
	return time.Duration(c.WorldStateTTL) * time.Second
}

// GetMarketDataTTL returns the market data cache duration
func (c *CacheConfig) GetMarketDataTTL() time.Duration {
	return time.Duration(c.MarketDataTTL) * time.Second
}
