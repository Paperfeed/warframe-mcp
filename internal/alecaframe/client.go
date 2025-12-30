package alecaframe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://stats.alecaframe.com/api"
)

// Client is an Alecaframe API client
type Client struct {
	httpClient  *http.Client
	userHash    string
	publicToken string
}

// NewClient creates a new Alecaframe API client
func NewClient(userHash, publicToken string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userHash:    userHash,
		publicToken: publicToken,
	}
}

// RelicInventory represents the user's relic inventory
type RelicInventory struct {
	Relics []Relic `json:"relics"`
}

// Relic represents a single relic in the inventory
type Relic struct {
	Name      string `json:"name"`
	Era       string `json:"era"`
	Tier      string `json:"tier"`
	Count     int    `json:"count"`
	IsVaulted bool   `json:"isVaulted"`
	Intact    int    `json:"intact"`
	Radiant   int    `json:"radiant"`
	Flawless  int    `json:"flawless"`
	Exceptional int  `json:"exceptional"`
}

// UserStats represents user trading and account statistics
type UserStats struct {
	TotalTrades   int     `json:"totalTrades"`
	PlatinumEarned int    `json:"platinumEarned"`
	PlatinumSpent  int    `json:"platinumSpent"`
	// Add more fields based on actual API response
}

// GetRelicInventory fetches the user's relic inventory
func (c *Client) GetRelicInventory() (*RelicInventory, error) {
	url := fmt.Sprintf("%s/relics/%s", baseURL, c.userHash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Add authorization header if public token is provided
	if c.publicToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.publicToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var inventory RelicInventory
	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &inventory, nil
}

// GetUserStats fetches the user's trading and account statistics
func (c *Client) GetUserStats() (*UserStats, error) {
	url := fmt.Sprintf("%s/stats/%s", baseURL, c.userHash)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var stats UserStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &stats, nil
}
