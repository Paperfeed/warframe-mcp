package market

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	baseURL = "https://api.warframe.market/v1"
)

// Client is a Warframe.market API client with rate limiting
type Client struct {
	httpClient  *http.Client
	rateLimiter *rateLimiter
	cache       *marketCache
}

// rateLimiter implements token bucket rate limiting (3 req/sec)
type rateLimiter struct {
	mu     sync.Mutex
	tokens int
	max    int
	refill time.Time
}

// marketCache holds cached market data
type marketCache struct {
	mu   sync.RWMutex
	data map[string]*cacheEntry
	ttl  time.Duration
}

type cacheEntry struct {
	data      interface{}
	timestamp time.Time
}

// NewClient creates a new Warframe.market API client
func NewClient(cacheTTL time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		rateLimiter: &rateLimiter{
			tokens: 3,
			max:    3,
			refill: time.Now(),
		},
		cache: &marketCache{
			data: make(map[string]*cacheEntry),
			ttl:  cacheTTL,
		},
	}
}

// ItemOrdersResponse represents the orders response
type ItemOrdersResponse struct {
	Payload struct {
		Orders []Order `json:"orders"`
	} `json:"payload"`
}

// Order represents a buy/sell order
type Order struct {
	ID           string    `json:"id"`
	User         OrderUser `json:"user"`
	Platinum     int       `json:"platinum"`
	Quantity     int       `json:"quantity"`
	OrderType    string    `json:"order_type"` // "buy" or "sell"
	Platform     string    `json:"platform"`
	CreationDate time.Time `json:"creation_date"`
	LastUpdate   time.Time `json:"last_update"`
	Visible      bool      `json:"visible"`
	ModRank      int       `json:"mod_rank,omitempty"`
}

// OrderUser represents the user who placed an order
type OrderUser struct {
	InGameName string `json:"ingame_name"`
	Status     string `json:"status"` // "ingame", "online", "offline"
	Reputation int    `json:"reputation"`
}

// ItemsResponse represents the items list response
type ItemsResponse struct {
	Payload struct {
		Items []Item `json:"items"`
	} `json:"payload"`
}

// Item represents a tradeable item
type Item struct {
	ID       string   `json:"id"`
	URLName  string   `json:"url_name"`
	Thumb    string   `json:"thumb"`
	ItemName string   `json:"item_name"`
	Tags     []string `json:"tags"`
}

// ItemStatistics represents price statistics
type ItemStatistics struct {
	Statistics48h  []StatPoint `json:"48hours"`
	Statistics90d  []StatPoint `json:"90days"`
	StatisticsClosed []StatPoint `json:"closed"`
}

// StatPoint represents a statistics data point
type StatPoint struct {
	DateTime    time.Time `json:"datetime"`
	Volume      int       `json:"volume"`
	MinPrice    int       `json:"min_price"`
	MaxPrice    int       `json:"max_price"`
	AvgPrice    float64   `json:"avg_price"`
	MedianPrice float64   `json:"median"`
	MovingAvg   float64   `json:"moving_avg"`
}

// allow implements token bucket rate limiting
func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Refill tokens every second
	if now.Sub(rl.refill) >= time.Second {
		rl.tokens = rl.max
		rl.refill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// wait blocks until a token is available
func (rl *rateLimiter) wait() {
	for !rl.allow() {
		time.Sleep(100 * time.Millisecond)
	}
}

// doRequest performs an HTTP request with rate limiting
func (c *Client) doRequest(endpoint string) ([]byte, error) {
	c.rateLimiter.wait()

	url := baseURL + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set required headers
	req.Header.Set("Platform", "pc")
	req.Header.Set("Language", "en")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// GetItemOrders fetches buy and sell orders for an item
func (c *Client) GetItemOrders(itemURLName string) ([]Order, error) {
	cacheKey := "orders:" + itemURLName

	// Check cache
	c.cache.mu.RLock()
	if entry, ok := c.cache.data[cacheKey]; ok {
		if time.Since(entry.timestamp) < c.cache.ttl {
			c.cache.mu.RUnlock()
			return entry.data.([]Order), nil
		}
	}
	c.cache.mu.RUnlock()

	endpoint := fmt.Sprintf("/items/%s/orders", itemURLName)
	data, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var response ItemOrdersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	orders := response.Payload.Orders

	// Update cache
	c.cache.mu.Lock()
	c.cache.data[cacheKey] = &cacheEntry{
		data:      orders,
		timestamp: time.Now(),
	}
	c.cache.mu.Unlock()

	return orders, nil
}

// SearchItems searches for items by name
func (c *Client) SearchItems(query string) ([]Item, error) {
	endpoint := "/items"
	data, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var response ItemsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Filter items by query
	query = strings.ToLower(query)
	var results []Item
	for _, item := range response.Payload.Items {
		if strings.Contains(strings.ToLower(item.ItemName), query) {
			results = append(results, item)
		}
	}

	return results, nil
}

// GetItemPrice gets the current average price for an item
func (c *Client) GetItemPrice(itemURLName string) (*PriceInfo, error) {
	orders, err := c.GetItemOrders(itemURLName)
	if err != nil {
		return nil, err
	}

	var sellOrders, buyOrders []int
	for _, order := range orders {
		if !order.Visible {
			continue
		}

		// Only count online/ingame users
		if order.User.Status == "ingame" || order.User.Status == "online" {
			if order.OrderType == "sell" {
				sellOrders = append(sellOrders, order.Platinum)
			} else if order.OrderType == "buy" {
				buyOrders = append(buyOrders, order.Platinum)
			}
		}
	}

	return &PriceInfo{
		ItemURLName:   itemURLName,
		SellOrders:    sellOrders,
		BuyOrders:     buyOrders,
		AvgSellPrice:  avg(sellOrders),
		AvgBuyPrice:   avg(buyOrders),
		MinSellPrice:  min(sellOrders),
		MaxBuyPrice:   max(buyOrders),
		TotalSellOrders: len(sellOrders),
		TotalBuyOrders:  len(buyOrders),
	}, nil
}

// PriceInfo holds aggregated price information
type PriceInfo struct {
	ItemURLName      string
	SellOrders       []int
	BuyOrders        []int
	AvgSellPrice     float64
	AvgBuyPrice      float64
	MinSellPrice     int
	MaxBuyPrice      int
	TotalSellOrders  int
	TotalBuyOrders   int
}

// URLEncode converts an item name to URL format
func URLEncode(itemName string) string {
	itemName = strings.ToLower(itemName)
	itemName = strings.ReplaceAll(itemName, " ", "_")
	return url.PathEscape(itemName)
}

// Helper functions
func avg(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums))
}

func min(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums {
		if n < m {
			m = n
		}
	}
	return m
}

func max(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums {
		if n > m {
			m = n
		}
	}
	return m
}
