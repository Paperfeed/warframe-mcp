# Warframe MCP Server

A Model Context Protocol (MCP) server for Warframe that provides intelligent recommendations based on your account data, current world state, and market conditions.

## Features

- **Account Integration**: Connect your Warframe account via Alecaframe API to access your inventory, relics, and progress
- **World State Awareness**: Real-time data on current events, rotations, vault status, and invasions
- **Market Intelligence**: Integration with Warframe.market for trading insights and pricing
- **Smart Recommendations**: Answer questions like:
  - "Which prime parts should I farm now?"
  - "What standing should I focus on?"
  - "Best way to gain credits based on current events?"
  - "What relics are most valuable right now?"
  - "Is Baro here and what's he selling?"
  - "What's the market price for [item]?"

## Available Tools

The MCP server provides 10 powerful tools:

### Recommendations
1. **recommend_prime_farming** - Get personalized prime part farming recommendations
2. **recommend_standing_farming** - Find the best syndicate standing opportunities
3. **recommend_credit_farming** - Optimize your credit farming strategy

### World State & Events
4. **get_world_state** - View current game state, events, and cycles
5. **get_baro_status** - Check Baro Ki'Teer's location and inventory
6. **get_varzia_status** - View Prime Resurgence rotation and schedule

### Market Data
7. **get_market_price** - Look up item prices on warframe.market
8. **search_market_items** - Search for tradeable items

### Your Account (Requires Alecaframe)
9. **get_my_relics** - View your personal relic inventory with counts and refinement levels
10. **get_my_stats** - View your trading statistics and account metrics

See [USAGE.md](USAGE.md) for detailed documentation and examples.

## Setup

### Prerequisites

- Go 1.21 or higher
- (Optional) Warframe account with Alecaframe integration for personalized relic recommendations

### Configuration

The server supports three configuration methods (in order of precedence):

**Option 1: Environment Variables (Recommended for MCP)**

Set environment variables in your MCP server configuration:
- `WARFRAME_USER_HASH` - Your Alecaframe user hash
- `WARFRAME_PUBLIC_TOKEN` - Your Alecaframe public token

This is the easiest method when using with Claude Desktop or other MCP clients.

**Option 2: Config File**

Create a `config.json` in one of these locations:
- Current directory: `./config.json`
- Home directory: `~/.config/warframe-mcp/config.json`
- System: `/etc/warframe-mcp/config.json`

```json
{
  "alecaframe": {
    "userHash": "your-user-hash-here",
    "publicToken": "your-public-token-here"
  },
  "cache": {
    "worldStateTTL": 300,
    "marketDataTTL": 600
  }
}
```

See `config.example.json` for a template.

**Option 3: No Configuration**

The server will work without any configuration, but Alecaframe features (personalized relic inventory) will be disabled.

### Installation

**On macOS/Linux:**
```bash
go build -o warframe-mcp
```

**On Windows (PowerShell or Command Prompt):**
```powershell
go build -o warframe-mcp.exe
```

> 📝 **Windows Users:** See [WINDOWS_SETUP.md](WINDOWS_SETUP.md) for a complete step-by-step Windows setup guide with troubleshooting tips.

### Usage with Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS, `%APPDATA%\Claude\claude_desktop_config.json` on Windows):

**Option 1: Basic Setup (No Alecaframe)**

macOS/Linux:
```json
{
  "mcpServers": {
    "warframe": {
      "command": "/absolute/path/to/warframe-mcp"
    }
  }
}
```

Windows:
```json
{
  "mcpServers": {
    "warframe": {
      "command": "C:\\Users\\YourUsername\\path\\to\\warframe-mcp.exe"
    }
  }
}
```

**Option 2: With Alecaframe Credentials (Recommended)**

macOS/Linux:
```json
{
  "mcpServers": {
    "warframe": {
      "command": "/absolute/path/to/warframe-mcp",
      "env": {
        "WARFRAME_USER_HASH": "your-user-hash-here",
        "WARFRAME_PUBLIC_TOKEN": "your-public-token-here"
      }
    }
  }
}
```

Windows:
```json
{
  "mcpServers": {
    "warframe": {
      "command": "C:\\Users\\YourUsername\\path\\to\\warframe-mcp.exe",
      "env": {
        "WARFRAME_USER_HASH": "your-user-hash-here",
        "WARFRAME_PUBLIC_TOKEN": "your-public-token-here"
      }
    }
  }
}
```

To get your credentials:
1. Open Warframe and go to the Stats tab
2. Generate a public token with 'relic' access enabled
3. Copy your userHash (keep this private!)
4. Add them to the `env` section above

Restart Claude Desktop, and you'll see the Warframe tools available in conversations!

### Quick Start

Once configured, try asking Claude:
- "What should I farm in Warframe right now?"
- "Is Baro here and what's worth buying?"
- "What's the best way to get credits today?"
- "Show me the current world state"

See [USAGE.md](USAGE.md) for more examples and detailed usage instructions.

## How It Works

The Warframe MCP server integrates multiple Warframe community APIs to provide comprehensive, context-aware recommendations:

### API Integrations

- **[Alecaframe API](https://docs.alecaframe.com/api)** - User account data and relic inventory
  - Provides your personal relic collection
  - Enables personalized farming recommendations

- **[WarframeStat.us API](https://docs.warframestat.us/)** - World state and event tracking
  - Real-time game state (alerts, invasions, fissures, sorties)
  - Baro Ki'Teer and Varzia schedules
  - Open world cycles (Cetus, Orb Vallis, Cambion Drift)

- **[Warframe.market API](https://warframe.market/api_docs)** - Trading and market data
  - Live buy/sell orders
  - Price trends and market statistics
  - Item search and discovery

### Intelligent Recommendations

The server combines data from all sources to provide smart, actionable advice:
- Cross-references your relics with active fissures
- Identifies high-value items from current vault rotations
- Suggests optimal activities based on time-limited events
- Provides market insights for trading decisions

### Caching & Rate Limiting

- Implements intelligent caching (5min for world state, 10min for market data)
- Respects API rate limits (3 req/sec for warframe.market)
- Minimizes API calls while keeping data fresh

## Project Structure

```
warframe-mcp/
├── main.go                          # MCP server entry point
├── internal/
│   ├── config/                      # Configuration management
│   ├── alecaframe/                  # Alecaframe API client
│   ├── worldstate/                  # WarframeStat.us API client
│   ├── market/                      # Warframe.market API client
│   └── recommender/                 # Recommendation engine
├── config.example.json              # Configuration template
└── README.md                        # This file
```

## Contributing

Contributions are welcome! Feel free to:
- Report bugs or suggest features via GitHub Issues
- Submit pull requests with improvements
- Add support for additional Warframe APIs
- Improve recommendation algorithms

## License

MIT

## Acknowledgments

Built using:
- [mcp-go](https://github.com/mark3labs/mcp-go) by mark3labs
- Warframe community APIs and the amazing WFCD team
- Claude by Anthropic for the MCP protocol
