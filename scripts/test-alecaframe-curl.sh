#!/bin/bash

# Script to test Alecaframe API endpoints using curl
# Usage: ./scripts/test-alecaframe-curl.sh YOUR_USER_HASH YOUR_TOKEN

set -e

USER_HASH="${1:-$ALECAFRAME_USER_HASH}"
TOKEN="${2:-$ALECAFRAME_PUBLIC_TOKEN}"

if [ -z "$USER_HASH" ] || [ -z "$TOKEN" ]; then
    echo "Error: Missing credentials"
    echo ""
    echo "Usage:"
    echo "  ./scripts/test-alecaframe-curl.sh YOUR_HASH YOUR_TOKEN"
    echo ""
    echo "Or set environment variables:"
    echo "  export ALECAFRAME_USER_HASH=your-hash"
    echo "  export ALECAFRAME_PUBLIC_TOKEN=your-token"
    echo "  ./scripts/test-alecaframe-curl.sh"
    exit 1
fi

echo "=== Alecaframe API Curl Test ==="
echo "User Hash: ${USER_HASH:0:8}***"
echo "Token: ${TOKEN:0:8}***"
echo ""

# Array of base URLs to try
BASE_URLS=(
    "https://stats.alecaframe.com/api"
    "https://stats.alecaframe.com/api/v1"
    "https://api.alecaframe.com"
    "https://api.alecaframe.com/v1"
)

# Endpoint patterns for relics
RELIC_PATTERNS=(
    "/relics/$USER_HASH"
    "/relic/$USER_HASH"
    "/users/$USER_HASH/relics"
    "/user/$USER_HASH/relics"
    "/$USER_HASH/relics"
)

# Endpoint patterns for stats
STATS_PATTERNS=(
    "/stats/$USER_HASH"
    "/statistics/$USER_HASH"
    "/users/$USER_HASH/stats"
    "/user/$USER_HASH/stats"
    "/$USER_HASH/stats"
)

test_url() {
    local url=$1
    local name=$2

    echo "Testing: $url"

    # Try Bearer token
    response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $TOKEN" "$url" 2>/dev/null || echo -e "\n000")
    status_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$status_code" = "200" ]; then
        echo "  ✓ SUCCESS with Bearer token!"
        echo "  Response preview:"
        echo "$body" | head -c 300
        echo ""
        echo "  Working URL: $url"
        echo "  Auth: Bearer token in Authorization header"
        echo ""
        return 0
    elif [ "$status_code" != "404" ] && [ "$status_code" != "401" ]; then
        echo "  Status: $status_code (unexpected, check response)"
        echo "  Response: ${body:0:100}"
    fi

    # Try X-API-Key header
    response=$(curl -s -w "\n%{http_code}" -H "X-API-Key: $TOKEN" "$url" 2>/dev/null || echo -e "\n000")
    status_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$status_code" = "200" ]; then
        echo "  ✓ SUCCESS with X-API-Key!"
        echo "  Response preview:"
        echo "$body" | head -c 300
        echo ""
        echo "  Working URL: $url"
        echo "  Auth: Token in X-API-Key header"
        echo ""
        return 0
    fi

    # Try as query parameter
    response=$(curl -s -w "\n%{http_code}" "$url?token=$TOKEN" 2>/dev/null || echo -e "\n000")
    status_code=$(echo "$response" | tail -n 1)

    if [ "$status_code" = "200" ]; then
        echo "  ✓ SUCCESS with query parameter!"
        echo "  Working URL: $url?token=..."
        echo "  Auth: Token as query parameter"
        echo ""
        return 0
    fi

    echo ""
    return 1
}

echo "=== Testing Relics Endpoint ==="
echo ""

found=false
for base in "${BASE_URLS[@]}"; do
    for pattern in "${RELIC_PATTERNS[@]}"; do
        url="$base$pattern"
        if test_url "$url" "relics"; then
            found=true
            break 2
        fi
    done
done

if [ "$found" = false ]; then
    echo "❌ No working relics endpoint found"
fi

echo ""
echo "=== Testing Stats Endpoint ==="
echo ""

found=false
for base in "${BASE_URLS[@]}"; do
    for pattern in "${STATS_PATTERNS[@]}"; do
        url="$base$pattern"
        if test_url "$url" "stats"; then
            found=true
            break 2
        fi
    done
done

if [ "$found" = false ]; then
    echo "❌ No working stats endpoint found"
fi

echo ""
echo "=== Test Complete ==="
echo ""
echo "If you found working endpoints, update internal/alecaframe/client.go"
echo "See TESTING_ALECAFRAME.md for more details"
