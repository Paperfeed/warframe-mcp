# Testing Alecaframe API Integration

This guide helps you test the Alecaframe API locally to verify your credentials.

## Correct API Endpoints (Updated)

Based on the Alecaframe Swagger documentation at https://stats.alecaframe.com/api/swagger/index.html:

1. **Get User Stats**: `/api/stats/{userHash}?secretToken={token}`
   ```bash
   curl -H 'accept: application/json' \
     'https://stats.alecaframe.com/api/stats/YOUR_HASH?secretToken=YOUR_TOKEN'
   ```

2. **Get Relic Inventory**: `/api/stats/public/getRelicInventory?publicToken={token}`
   ```bash
   curl -H 'accept: application/json' \
     'https://stats.alecaframe.com/api/stats/public/getRelicInventory?publicToken=YOUR_TOKEN'
   ```

Note: Relic inventory doesn't require the userHash, only the public token!

## Prerequisites

1. **Get Your Alecaframe Credentials**
   - Go to [Alecaframe](https://alecaframe.com) and link your Warframe account if you haven't already
   - Navigate to the "Stats" tab on Alecaframe
   - Click on the "Data export & API" button
   - Click "Copy my API token" - this is your `ALECAFRAME_USER_HASH` (keep this private!)
   - For `ALECAFRAME_PUBLIC_TOKEN`: Location currently unknown - check Alecaframe documentation or contact support

## Method 1: Quick Test with curl (Recommended)

Test directly using the correct endpoints:

```bash
# Test relic inventory (doesn't need userHash)
curl -v -H 'accept: application/json' \
  'https://stats.alecaframe.com/api/stats/public/getRelicInventory?publicToken=YOUR_TOKEN'

# Test user stats (needs both userHash and token)
curl -v -H 'accept: application/json' \
  'https://stats.alecaframe.com/api/stats/YOUR_HASH?secretToken=YOUR_TOKEN'
```

Expected responses:
- **200 OK**: Success! You'll see your data
- **401 Unauthorized**: Wrong or expired token
- **404 Not Found**: Wrong userHash

## Method 2: Manual Testing with cURL

If you know your credentials, you can test manually:

### Test 1: Relics Endpoint

```bash
# Try pattern 1: /api/relics/{hash}
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://stats.alecaframe.com/api/relics/YOUR_HASH

# Try pattern 2: /api/v1/relics/{hash}
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://stats.alecaframe.com/api/v1/relics/YOUR_HASH

# Try pattern 3: /api/users/{hash}/relics
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://stats.alecaframe.com/api/users/YOUR_HASH/relics

# Try with X-API-Key header instead
curl -H "X-API-Key: YOUR_TOKEN" \
  https://stats.alecaframe.com/api/relics/YOUR_HASH
```

### Test 2: Stats Endpoint

```bash
# Try pattern 1: /api/stats/{hash}
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://stats.alecaframe.com/api/stats/YOUR_HASH

# Try pattern 2: /api/v1/stats/{hash}
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://stats.alecaframe.com/api/v1/stats/YOUR_HASH
```

### Check Response Codes

- **200 OK**: Success! The endpoint works
- **401 Unauthorized**: Authentication issue (wrong token or auth method)
- **403 Forbidden**: Token doesn't have required permissions
- **404 Not Found**: Wrong URL pattern

## Method 3: Using the MCP Server with Debug Logging

Add debug logging to see exactly what's being called:

```bash
# Set the environment variables
export ALECAFRAME_USER_HASH="your-hash"
export ALECAFRAME_PUBLIC_TOKEN="your-token"

# Run the MCP server
./warframe-mcp
```

Then in Claude Desktop, try:
- "Show me my relics"
- "Show me my stats"

Check the terminal output for error messages and URLs being called.

## Common Issues

### Issue 1: 404 Not Found
**Cause**: Wrong endpoint URL pattern
**Solution**: Use the test utility to find the correct URL pattern

### Issue 2: 401 Unauthorized
**Cause**: Wrong authentication method or invalid token
**Solution**:
- Verify your token is correct
- Check if it has 'relic' access enabled
- Try different auth methods (Bearer vs API Key)

### Issue 3: 403 Forbidden
**Cause**: Token doesn't have required permissions
**Solution**: Regenerate your token in Warframe with proper permissions

### Issue 4: Token Expired
**Cause**: Tokens expire after 1 year
**Solution**: Generate a new token in Warframe

## Updating the Code

Once you find the working URL pattern and auth method:

1. **Update `internal/alecaframe/client.go`:**

```go
// Update the endpoints
func (c *Client) GetRelicInventory() (*RelicInventory, error) {
    // Change this line to match what works:
    url := fmt.Sprintf("%s/CORRECT_PATTERN/%s", baseURL, c.userHash)

    // Update auth method if needed:
    req.Header.Set("X-API-Key", c.publicToken)  // or whatever works
    ...
}
```

2. **Test again:**

```bash
go test ./internal/alecaframe -v
```

3. **Rebuild:**

```bash
make build
```

## Getting Help

If you're still having issues:

1. **Check the Alecaframe documentation**: https://docs.alecaframe.com/api
2. **Verify your credentials are correct** in the Warframe Stats tab
3. **Check if Alecaframe API is operational**
4. **Open an issue** with the output from the test utility

## Example Output

Successful test output should look like:

```
Testing: https://stats.alecaframe.com/api/v1/relics/abc123def456
  [Bearer token in Authorization header] Status: 200
  ✓ SUCCESS! Response preview:
    {
      "relics": [
        {
          "name": "A5",
          "era": "Axi",
          "count": 10,
          ...
        }
      ]
    }

  Working URL: https://stats.alecaframe.com/api/v1/relics/abc123def456
  Auth method: Bearer token in Authorization header
```
