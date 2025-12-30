# Warframe MCP - Windows Setup Guide

Complete guide for setting up and running the Warframe MCP server on Windows.

## Prerequisites

1. **Go (Golang)** - Download and install from [https://go.dev/dl/](https://go.dev/dl/)
   - Choose the Windows installer (.msi file)
   - Follow the installation wizard
   - Verify installation by opening PowerShell and running: `go version`

2. **Git (Optional)** - Only needed if cloning from repository
   - Download from [https://git-scm.com/download/win](https://git-scm.com/download/win)

3. **Claude Desktop** - Download from [https://claude.ai/download](https://claude.ai/download)

## Step-by-Step Installation

### 1. Get the Code

**Option A: Clone from Git**
```powershell
git clone https://github.com/paperfeed/warframe-mcp.git
cd warframe-mcp
```

**Option B: Download ZIP**
- Download the repository as a ZIP file
- Extract to a location like `C:\warframe-mcp`
- Open PowerShell and navigate to the folder:
```powershell
cd C:\warframe-mcp
```

### 2. Build the Application

In PowerShell, run:
```powershell
go build -o warframe-mcp.exe
```

This creates `warframe-mcp.exe` in the current directory.

### 3. (Optional) Get Your Alecaframe Credentials

For personalized recommendations based on your relic inventory:

1. Go to [Alecaframe](https://alecaframe.com) and link your Warframe account if you haven't already
2. Navigate to the "Stats" tab on Alecaframe
3. Click on the "Data export & API" button
4. Click "Copy my API token" - this is your `ALECAFRAME_USER_HASH` (keep this private!)
5. For `ALECAFRAME_PUBLIC_TOKEN`: Location currently unknown - check Alecaframe documentation or contact support
6. Save these for the next step

### 4. Configure Claude Desktop

#### Find the Config File

The Claude Desktop configuration file is located at:
```
%APPDATA%\Claude\claude_desktop_config.json
```

To open it:
1. Press `Win + R`
2. Type `%APPDATA%\Claude` and press Enter
3. Open `claude_desktop_config.json` with Notepad or your preferred editor
   - If the file doesn't exist, create it

#### Basic Configuration (Without Alecaframe)

```json
{
  "mcpServers": {
    "warframe": {
      "command": "C:\\warframe-mcp\\warframe-mcp.exe"
    }
  }
}
```

**Note:** Replace `C:\\warframe-mcp\\warframe-mcp.exe` with the actual path where you built the executable. Use double backslashes (`\\`) in the path!

#### Full Configuration (With Alecaframe - Recommended)

```json
{
  "mcpServers": {
    "warframe": {
      "command": "C:\\warframe-mcp\\warframe-mcp.exe",
      "env": {
        "ALECAFRAME_USER_HASH": "your-user-hash-here",
        "ALECAFRAME_PUBLIC_TOKEN": "your-public-token-here"
      }
    }
  }
}
```

### 5. Restart Claude Desktop

Close Claude Desktop completely and restart it. The Warframe tools will now be available!

## Testing the Installation

After configuring Claude Desktop, try asking:
- "What should I farm in Warframe right now?"
- "Is Baro here?"
- "Show me the current world state"

You should see Claude using the Warframe tools to answer your questions.

## Troubleshooting

### "Go is not recognized as an internal or external command"

- Go is not installed or not in your PATH
- Reinstall Go and make sure to check "Add to PATH" during installation
- Restart PowerShell after installation

### Claude Desktop doesn't show Warframe tools

1. Check that the config file path is correct
2. Verify the `.exe` path uses double backslashes (`\\`)
3. Make sure the path is absolute (starts with drive letter like `C:\\`)
4. Check Claude Desktop logs in `%APPDATA%\Claude\logs`
5. Try restarting Claude Desktop completely

### "The system cannot find the path specified"

- The path to `warframe-mcp.exe` is incorrect
- Use the full absolute path, for example:
  - ✅ `C:\\Users\\YourName\\warframe-mcp\\warframe-mcp.exe`
  - ❌ `warframe-mcp.exe` (relative path won't work)
  - ❌ `C:\Users\YourName\warframe-mcp\warframe-mcp.exe` (single backslash won't work in JSON)

### Environment variables not working

If you're having issues with `env` in the config, you can alternatively create a `config.json` file:

1. Create `config.json` in the same directory as `warframe-mcp.exe`:
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

2. Update your Claude Desktop config to just:
```json
{
  "mcpServers": {
    "warframe": {
      "command": "C:\\warframe-mcp\\warframe-mcp.exe"
    }
  }
}
```

## Quick Path Finding

To get the full path to your executable:

1. Open PowerShell
2. Navigate to the warframe-mcp directory:
   ```powershell
   cd C:\path\to\warframe-mcp
   ```
3. Run:
   ```powershell
   (Get-Item .\warframe-mcp.exe).FullName
   ```
4. This will print the full path - copy this and replace backslashes with double backslashes for the JSON config

## Alternative: Using WSL (Windows Subsystem for Linux)

If you have WSL installed, you can follow the Linux instructions instead:

1. Open WSL terminal
2. Clone and build:
   ```bash
   git clone https://github.com/paperfeed/warframe-mcp.git
   cd warframe-mcp
   go build -o warframe-mcp
   ```
3. In Claude Desktop config, use the WSL path:
   ```json
   {
     "mcpServers": {
       "warframe": {
         "command": "wsl",
         "args": ["-e", "/home/yourusername/warframe-mcp/warframe-mcp"]
       }
     }
   }
   ```

## Need More Help?

- Check the main [README.md](README.md) for general documentation
- See [USAGE.md](USAGE.md) for detailed tool usage and examples
- Report issues on GitHub
