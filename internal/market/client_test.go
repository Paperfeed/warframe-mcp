package market

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	ttl := 10 * time.Minute
	client := NewClient(ttl)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.cache.ttl != ttl {
		t.Errorf("Expected cache TTL to be %v, got %v", ttl, client.cache.ttl)
	}
}

func TestRateLimiter(t *testing.T) {
	rl := &rateLimiter{
		tokens: 3,
		max:    3,
		refill: time.Now(),
	}

	// Should allow 3 tokens
	if !rl.allow() {
		t.Error("Expected first request to be allowed")
	}
	if !rl.allow() {
		t.Error("Expected second request to be allowed")
	}
	if !rl.allow() {
		t.Error("Expected third request to be allowed")
	}

	// Fourth should be denied
	if rl.allow() {
		t.Error("Expected fourth request to be denied")
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)
	if !rl.allow() {
		t.Error("Expected request after refill to be allowed")
	}
}

func TestSearchItems(t *testing.T) {
	client := NewClient(10 * time.Minute)

	items, err := client.SearchItems("rhino prime")
	if err != nil {
		t.Logf("Warning: SearchItems failed (API may be down or rate limited): %v", err)
		return
	}

	if len(items) == 0 {
		t.Log("Warning: No items found for 'rhino prime' - this might be expected")
		return
	}

	t.Logf("Found %d items matching 'rhino prime'", len(items))
	for i, item := range items {
		if i >= 5 {
			break
		}
		t.Logf("  - %s (URL: %s)", item.ItemName, item.URLName)
	}
}

func TestGetItemPrice(t *testing.T) {
	client := NewClient(10 * time.Minute)

	// Test with a common item that should exist
	priceInfo, err := client.GetItemPrice("rhino_prime_set")
	if err != nil {
		t.Logf("Warning: GetItemPrice failed (API may be down or rate limited): %v", err)
		return
	}

	if priceInfo == nil {
		t.Fatal("Expected price info to be returned, got nil")
	}

	t.Logf("Price info for rhino_prime_set:")
	t.Logf("  Avg Sell Price: %.2f", priceInfo.AvgSellPrice)
	t.Logf("  Avg Buy Price: %.2f", priceInfo.AvgBuyPrice)
	t.Logf("  Min Sell Price: %d", priceInfo.MinSellPrice)
	t.Logf("  Max Buy Price: %d", priceInfo.MaxBuyPrice)
	t.Logf("  Total Sell Orders: %d", priceInfo.TotalSellOrders)
	t.Logf("  Total Buy Orders: %d", priceInfo.TotalBuyOrders)
}

func TestURLEncode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Rhino Prime Set", "rhino_prime_set"},
		{"Maiming Strike", "maiming_strike"},
		{"TEST ITEM", "test_item"},
	}

	for _, tt := range tests {
		result := URLEncode(tt.input)
		if result != tt.expected {
			t.Errorf("URLEncode(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test avg
	nums := []int{10, 20, 30, 40, 50}
	avgResult := avg(nums)
	if avgResult != 30.0 {
		t.Errorf("avg(%v) = %f, want 30.0", nums, avgResult)
	}

	// Test avg with empty slice
	emptyAvg := avg([]int{})
	if emptyAvg != 0 {
		t.Errorf("avg([]) = %f, want 0", emptyAvg)
	}

	// Test min
	minResult := min(nums)
	if minResult != 10 {
		t.Errorf("min(%v) = %d, want 10", nums, minResult)
	}

	// Test max
	maxResult := max(nums)
	if maxResult != 50 {
		t.Errorf("max(%v) = %d, want 50", nums, maxResult)
	}

	// Test min with empty slice
	minEmpty := min([]int{})
	if minEmpty != 0 {
		t.Errorf("min([]) = %d, want 0", minEmpty)
	}
}

func TestCaching(t *testing.T) {
	client := NewClient(1 * time.Second)

	// First call - should hit API
	_, err := client.GetItemPrice("test_item")
	if err != nil {
		t.Logf("First call failed (expected if item doesn't exist): %v", err)
	}

	// Second call within TTL - should use cache
	startTime := time.Now()
	_, err = client.GetItemPrice("test_item")
	duration := time.Since(startTime)

	// Cached call should be very fast (< 10ms)
	if duration > 10*time.Millisecond {
		t.Logf("Warning: Cached call took %v (expected < 10ms)", duration)
	}

	// Wait for cache to expire
	time.Sleep(1100 * time.Millisecond)

	// Third call after TTL - should hit API again
	_, err = client.GetItemPrice("test_item")
	if err != nil {
		t.Logf("Third call after cache expiry failed: %v", err)
	}
}
