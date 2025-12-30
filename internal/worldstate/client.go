package worldstate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL = "https://api.warframestat.us/pc"
)

// Client is a WarframeStat.us API client
type Client struct {
	httpClient *http.Client
	cache      *cache
}

// cache holds cached world state data
type cache struct {
	mu         sync.RWMutex
	data       *WorldState
	lastUpdate time.Time
	ttl        time.Duration
}

// NewClient creates a new world state API client
func NewClient(cacheTTL time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: &cache{
			ttl: cacheTTL,
		},
	}
}

// WorldState represents the current game state
type WorldState struct {
	Timestamp      string          `json:"timestamp"`
	Alerts         []Alert         `json:"alerts"`
	Arbitration    *Arbitration    `json:"arbitration"`
	Sortie         *Sortie         `json:"sortie"`
	SyndicateMissions []SyndicateMission `json:"syndicateMissions"`
	Fissures       []Fissure       `json:"fissures"`
	Invasions      []Invasion      `json:"invasions"`
	VoidTrader     *VoidTrader     `json:"voidTrader"`
	VaultTrader    *VaultTrader    `json:"vaultTrader"`
	Events         []Event         `json:"events"`
	Nightwave      *Nightwave      `json:"nightwave"`
	CambionDrift   *CambionCycle   `json:"cambionCycle"`
	CetusCycle     *CetusCycle     `json:"cetusCycle"`
	VallisEarth    *VallisEarth    `json:"vallisCycle"`
	Zariman        *Zariman        `json:"zarimanCycle"`
}

// Alert represents an alert mission
type Alert struct {
	ID          string      `json:"id"`
	Activation  string      `json:"activation"`
	Expiry      string      `json:"expiry"`
	Mission     Mission     `json:"mission"`
	Active      bool        `json:"active"`
}

// Mission represents a mission
type Mission struct {
	Node         string   `json:"node"`
	Type         string   `json:"type"`
	Faction      string   `json:"faction"`
	Reward       *Reward  `json:"reward"`
	MinEnemyLevel int     `json:"minEnemyLevel"`
	MaxEnemyLevel int     `json:"maxEnemyLevel"`
}

// Reward represents mission rewards
type Reward struct {
	Credits      int      `json:"credits"`
	Items        []string `json:"items"`
	CountedItems []Item   `json:"countedItems"`
}

// Item represents a counted item
type Item struct {
	Count int    `json:"count"`
	Type  string `json:"type"`
}

// Arbitration represents the current arbitration
type Arbitration struct {
	Activation string  `json:"activation"`
	Expiry     string  `json:"expiry"`
	Node       string  `json:"node"`
	Enemy      string  `json:"enemy"`
	Type       string  `json:"type"`
	Active     bool    `json:"active"`
}

// Sortie represents the daily sortie
type Sortie struct {
	ID         string          `json:"id"`
	Activation string          `json:"activation"`
	Expiry     string          `json:"expiry"`
	Variants   []SortieVariant `json:"variants"`
	Boss       string          `json:"boss"`
	Faction    string          `json:"faction"`
	Active     bool            `json:"active"`
}

// SortieVariant represents a sortie mission variant
type SortieVariant struct {
	Node          string `json:"node"`
	MissionType   string `json:"missionType"`
	Modifier      string `json:"modifier"`
	ModifierDescription string `json:"modifierDescription"`
}

// SyndicateMission represents a syndicate mission
type SyndicateMission struct {
	Syndicate string   `json:"syndicate"`
	Nodes     []string `json:"nodes"`
	Jobs      []Job    `json:"jobs"`
}

// Job represents a syndicate job
type Job struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	StandingStages []int  `json:"standingStages"`
}

// Fissure represents a void fissure
type Fissure struct {
	ID         string `json:"id"`
	Node       string `json:"node"`
	Tier       string `json:"tier"`
	TierNum    int    `json:"tierNum"`
	MissionType string `json:"missionType"`
	Enemy      string `json:"enemy"`
	Activation string `json:"activation"`
	Expiry     string `json:"expiry"`
	Active     bool   `json:"active"`
	IsStorm    bool   `json:"isStorm"`
	IsHard     bool   `json:"isHard"`
}

// Invasion represents an invasion
type Invasion struct {
	ID          string      `json:"id"`
	Node        string      `json:"node"`
	Defender    Faction     `json:"defender"`
	Attacker    Faction     `json:"attacker"`
	Activation  string      `json:"activation"`
	Completed   bool        `json:"completed"`
	Completion  float64     `json:"completion"`
}

// Faction represents an invasion faction
type Faction struct {
	Faction string  `json:"faction"`
	Reward  *Reward `json:"reward"`
}

// VoidTrader represents Baro Ki'Teer
type VoidTrader struct {
	ID         string    `json:"id"`
	Activation string    `json:"activation"`
	Expiry     string    `json:"expiry"`
	Character  string    `json:"character"`
	Location   string    `json:"location"`
	Inventory  []VoidItem `json:"inventory"`
	Active     bool      `json:"active"`
}

// VoidItem represents a Baro item
type VoidItem struct {
	Item    string `json:"item"`
	Ducats  int    `json:"ducats"`
	Credits int    `json:"credits"`
}

// VaultTrader represents Varzia (Prime Resurgence)
type VaultTrader struct {
	ID         string      `json:"id"`
	Activation string      `json:"activation"`
	Expiry     string      `json:"expiry"`
	Character  string      `json:"character"`
	Location   string      `json:"location"`
	Schedule   []Schedule  `json:"schedule"`
	Active     bool        `json:"active"`
}

// Schedule represents vault trader schedule
type Schedule struct {
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Items      []string `json:"items"`
	Featured   []string `json:"featured"`
}

// Event represents a game event
type Event struct {
	ID          string   `json:"id"`
	Activation  string   `json:"activation"`
	Expiry      string   `json:"expiry"`
	Description string   `json:"description"`
	Tooltip     string   `json:"tooltip"`
	Active      bool     `json:"active"`
	Rewards     []Reward `json:"rewards"`
}

// Nightwave represents the current nightwave season
type Nightwave struct {
	ID         string `json:"id"`
	Activation string `json:"activation"`
	Expiry     string `json:"expiry"`
	Season     int    `json:"season"`
	Tag        string `json:"tag"`
	Phase      int    `json:"phase"`
	Active     bool   `json:"active"`
}

// CambionCycle represents Cambion Drift time cycle
type CambionCycle struct {
	State  string `json:"state"`
	Expiry string `json:"expiry"`
}

// CetusCycle represents Plains of Eidolon time cycle
type CetusCycle struct {
	IsDay      bool   `json:"isDay"`
	State      string `json:"state"`
	Expiry     string `json:"expiry"`
	TimeLeft   string `json:"timeLeft"`
}

// VallisEarth represents Orb Vallis time cycle
type VallisEarth struct {
	IsWarm     bool   `json:"isWarm"`
	State      string `json:"state"`
	Expiry     string `json:"expiry"`
	TimeLeft   string `json:"timeLeft"`
}

// Zariman represents Zariman time cycle
type Zariman struct {
	State      string `json:"state"`
	Expiry     string `json:"expiry"`
	TimeLeft   string `json:"timeLeft"`
}

// GetWorldState fetches the current world state with caching
func (c *Client) GetWorldState() (*WorldState, error) {
	c.cache.mu.RLock()
	if c.cache.data != nil && time.Since(c.cache.lastUpdate) < c.cache.ttl {
		defer c.cache.mu.RUnlock()
		return c.cache.data, nil
	}
	c.cache.mu.RUnlock()

	// Fetch fresh data
	resp, err := c.httpClient.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("fetching world state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var ws WorldState
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Update cache
	c.cache.mu.Lock()
	c.cache.data = &ws
	c.cache.lastUpdate = time.Now()
	c.cache.mu.Unlock()

	return &ws, nil
}

// GetVaultTrader fetches the current vault trader (Varzia) status
func (c *Client) GetVaultTrader() (*VaultTrader, error) {
	resp, err := c.httpClient.Get(baseURL + "/vaultTrader")
	if err != nil {
		return nil, fmt.Errorf("fetching vault trader: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var vt VaultTrader
	if err := json.NewDecoder(resp.Body).Decode(&vt); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &vt, nil
}
