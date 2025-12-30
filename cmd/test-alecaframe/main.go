package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	userHash := flag.String("hash", "", "Alecaframe user hash")
	publicToken := flag.String("token", "", "Alecaframe public token")
	endpoint := flag.String("endpoint", "relics", "Endpoint to test (relics or stats)")
	flag.Parse()

	// Try to get from environment if not provided
	if *userHash == "" {
		*userHash = os.Getenv("ALECAFRAME_USER_HASH")
	}
	if *publicToken == "" {
		*publicToken = os.Getenv("ALECAFRAME_PUBLIC_TOKEN")
	}

	if *userHash == "" || *publicToken == "" {
		fmt.Println("Error: Must provide user hash and public token")
		fmt.Println("\nUsage:")
		fmt.Println("  go run cmd/test-alecaframe/main.go -hash YOUR_HASH -token YOUR_TOKEN -endpoint relics")
		fmt.Println("\nOr set environment variables:")
		fmt.Println("  export ALECAFRAME_USER_HASH=your-hash")
		fmt.Println("  export ALECAFRAME_PUBLIC_TOKEN=your-token")
		os.Exit(1)
	}

	fmt.Println("=== Alecaframe API Test ===")
	fmt.Printf("User Hash: %s\n", maskString(*userHash))
	fmt.Printf("Token: %s\n", maskString(*publicToken))
	fmt.Printf("Testing endpoint: %s\n\n", *endpoint)

	// Test different URL patterns
	baseURLs := []string{
		"https://stats.alecaframe.com/api",
		"https://stats.alecaframe.com/api/v1",
		"https://api.alecaframe.com",
		"https://api.alecaframe.com/v1",
	}

	endpointPatterns := map[string][]string{
		"relics": {
			"/relics/%s",
			"/relic/%s",
			"/users/%s/relics",
			"/user/%s/relics",
			"/%s/relics",
		},
		"stats": {
			"/stats/%s",
			"/statistics/%s",
			"/users/%s/stats",
			"/user/%s/stats",
			"/%s/stats",
		},
	}

	patterns, ok := endpointPatterns[*endpoint]
	if !ok {
		fmt.Printf("Unknown endpoint: %s\n", *endpoint)
		os.Exit(1)
	}

	fmt.Println("Testing different URL patterns...\n")

	for _, base := range baseURLs {
		for _, pattern := range patterns {
			url := base + fmt.Sprintf(pattern, *userHash)
			testURL(url, *publicToken)
		}
	}
}

func testURL(url, token string) {
	fmt.Printf("Testing: %s\n", url)

	client := &http.Client{Timeout: 10 * time.Second}

	// Try different authentication methods
	authMethods := []struct {
		name   string
		modify func(*http.Request)
	}{
		{
			name: "Bearer token in Authorization header",
			modify: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer "+token)
			},
		},
		{
			name: "Token in Authorization header (no Bearer)",
			modify: func(req *http.Request) {
				req.Header.Set("Authorization", token)
			},
		},
		{
			name: "Token in X-API-Key header",
			modify: func(req *http.Request) {
				req.Header.Set("X-API-Key", token)
			},
		},
		{
			name: "Token as query parameter",
			modify: func(req *http.Request) {
				q := req.URL.Query()
				q.Set("token", token)
				req.URL.RawQuery = q.Encode()
			},
		},
	}

	for _, auth := range authMethods {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("  [%s] Error creating request: %v\n", auth.name, err)
			continue
		}

		auth.modify(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("  [%s] Error: %v\n", auth.name, err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("  [%s] Status: %d\n", auth.name, resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("  ✓ SUCCESS! Response preview:\n")
			var prettyJSON map[string]interface{}
			if err := json.Unmarshal(body, &prettyJSON); err == nil {
				formatted, _ := json.MarshalIndent(prettyJSON, "    ", "  ")
				preview := string(formatted)
				if len(preview) > 500 {
					preview = preview[:500] + "..."
				}
				fmt.Printf("    %s\n", preview)
			} else {
				preview := string(body)
				if len(preview) > 200 {
					preview = preview[:200] + "..."
				}
				fmt.Printf("    %s\n", preview)
			}
			fmt.Printf("\n  Working URL: %s\n", url)
			fmt.Printf("  Auth method: %s\n\n", auth.name)
			return // Found working combination
		} else if resp.StatusCode != 404 && resp.StatusCode != 401 && resp.StatusCode != 403 {
			fmt.Printf("    Response: %s\n", string(body)[:min(200, len(body))])
		}
	}

	fmt.Println()
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
