package alecaframe

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Type       string `json:"type"`       // Lith, Meso, Neo, Axi, Requiem
	Refinement string `json:"refinement"` // Intact, Exceptional, Flawless, Radiant
	Name       string `json:"name"`       // e.g., "L1", "B21"
	Count      uint32 `json:"count"`      // Number of this specific relic
}

// UserStats represents user trading and account statistics
type UserStats struct {
	TotalTrades   int     `json:"totalTrades"`
	PlatinumEarned int    `json:"platinumEarned"`
	PlatinumSpent  int    `json:"platinumSpent"`
	// Add more fields based on actual API response
}

// relicTypeToString converts a relic type byte to a human-readable string
func relicTypeToString(relicType uint8) string {
	switch relicType {
	case 0:
		return "Lith"
	case 1:
		return "Meso"
	case 2:
		return "Neo"
	case 3:
		return "Axi"
	case 4:
		return "Requiem"
	default:
		return fmt.Sprintf("Unknown(%d)", relicType)
	}
}

// refinementToString converts a refinement byte to a human-readable string
func refinementToString(refinement uint8) string {
	switch refinement {
	case 0:
		return "Intact"
	case 1, 4:
		return "Exceptional"
	case 2, 5:
		return "Flawless"
	case 3, 6:
		return "Radiant"
	default:
		return fmt.Sprintf("Unknown(%d)", refinement)
	}
}

// GetRelicInventory fetches the user's relic inventory
// The API returns binary data in the following format (little endian):
// - Uint32: Number of relics
// - For each relic (9 bytes):
//   - Uint8: Relic type (0=Lith, 1=Meso, 2=Neo, 3=Axi, 4=Requiem)
//   - Uint8: Relic refinement (0=Intact, 1=Exceptional, 2=Flawless, 3=Radiant)
//   - char[3]: Name (e.g., "L1", "B21")
//   - Uint32: Count
func (c *Client) GetRelicInventory() (*RelicInventory, error) {
	endpoint := fmt.Sprintf("%s/stats/public/getRelicInventory?publicToken=%s", baseURL, url.QueryEscape(c.publicToken))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// The API returns a JSON string containing base64 data
	var base64String string
	if err := json.NewDecoder(resp.Body).Decode(&base64String); err != nil {
		return nil, fmt.Errorf("decoding JSON response: %w", err)
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, fmt.Errorf("decoding base64 response: %w", err)
	}

	// Check minimum size (at least 4 bytes for the count)
	if len(data) < 4 {
		return nil, fmt.Errorf("response too short: got %d bytes, need at least 4", len(data))
	}

	// Read number of relics (little endian Uint32)
	numRelics := binary.LittleEndian.Uint32(data[0:4])

	// Validate data size
	expectedSize := 4 + (numRelics * 9)
	if uint32(len(data)) != expectedSize {
		return nil, fmt.Errorf("unexpected data size: got %d bytes, expected %d for %d relics", len(data), expectedSize, numRelics)
	}

	// Parse each relic (9 bytes each)
	relics := make([]Relic, 0, numRelics)
	offset := uint32(4)

	for i := uint32(0); i < numRelics; i++ {
		relicType := data[offset]
		refinement := data[offset+1]
		nameBytes := data[offset+2 : offset+5]
		count := binary.LittleEndian.Uint32(data[offset+5 : offset+9])

		// Convert name bytes to string (null-terminated)
		name := string(nameBytes)
		for idx := 0; idx < len(name); idx++ {
			if name[idx] == 0 {
				name = name[:idx]
				break
			}
		}

		relics = append(relics, Relic{
			Type:       relicTypeToString(relicType),
			Refinement: refinementToString(refinement),
			Name:       name,
			Count:      count,
		})

		offset += 9
	}

	return &RelicInventory{Relics: relics}, nil
}

// GetUserStats fetches the user's trading and account statistics
// NOTE: This endpoint requires a secretToken which may be different from publicToken.
// Currently using publicToken as secretToken - if you get 401 errors, you may need
// to obtain a separate secretToken from Alecaframe.
func (c *Client) GetUserStats() (*UserStats, error) {
	// Using publicToken as secretToken for now - this may need to be a separate credential
	endpoint := fmt.Sprintf("%s/stats/%s?secretToken=%s", baseURL, url.PathEscape(c.userHash), url.QueryEscape(c.publicToken))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

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
