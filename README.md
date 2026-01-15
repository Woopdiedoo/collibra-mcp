# Collibra MCP Server

A Model Context Protocol (MCP) server that provides AI agents with access to Collibra Data Governance Center capabilities including data asset discovery, business glossary queries, and detailed asset information retrieval.

> **Note:** This is a fork of the original [Collibra MCP Server](https://github.com/collibra/chip) by [Collibra](https://www.collibra.com/), adapted for PGGM's SSO authentication requirements.

## Overview

This Go-based MCP server acts as a bridge between AI applications and Collibra, enabling intelligent data discovery and governance operations through the following tools:

- [`asset_details_get`](pkg/tools/get_asset_details.go) - Retrieve detailed information about specific assets by UUID
- [`asset_keyword_search`](pkg/tools/keyword_search.go) - Wildcard keyword search for assets
- [`asset_types_list`](pkg/tools/list_asset_types.go) - List available asset types
- [`business_glossary_discover`](pkg/tools/ask_glossary.go) - Ask questions about terms and definitions
- [`data_classification_match_add`](pkg/tools/add_data_classification_match.go) - Associate a data class with an asset
- [`data_classification_match_remove`](pkg/tools/remove_data_classification_match.go) - Remove a classification match
- [`data_classification_match_search`](pkg/tools/find_data_classification_matches.go) - Find associations between data classes and assets
- [`data_assets_discover`](pkg/tools/ask_dad.go) - Query available data assets using natural language
- [`data_class_search`](pkg/tools/search_data_classes.go) - Search for data classes with filters
- [`data_contract_list`](pkg/tools/list_data_contracts.go) - List data contracts with pagination
- [`data_contract_manifest_pull`](pkg/tools/pull_data_contract_manifest.go) - Download manifest for a data contract
- [`data_contract_manifest_push`](pkg/tools/push_data_contract_manifest.go) - Upload manifest for a data contract

## Quick Start

### Prerequisites

- **Go 1.21+** - [Download Go](https://go.dev/dl/)
- Access to a Collibra Data Governance Center instance
- SSO credentials (Azure AD / SAML)

### Installation

Build from source:
```bash
git clone https://github.com/Woopdiedoo/collibra-mcp.git
cd collibra-mcp
go mod download
go build -o chip ./cmd/chip
```

## Running and Configuration

### Authentication

The server uses cookie-based authentication with SSO support.

> **Why cookie-based instead of full SSO automation?**  
> Collibra uses SAML-based SSO, which is designed for browser-based flows and lacks a headless/programmatic grant type like OAuth2's client credentials. Automating SAML would require simulating a browser to handle redirects and form posts. Additionally, corporate environments with Azure AD and Intune device management only trust enrolled browser profiles, blocking automated browser sessions. By using your existing browser and manually copying the session cookie, we work around both limitations.

```bash
# First-time setup: authenticate via browser
./chip --api-url "https://pggm.collibra.com" --sso-auth
```

This will:
1. Open your default browser to the Collibra login page
2. Allow you to complete SSO authentication
3. Prompt you to paste the session cookie from browser DevTools (see below)
4. Cache the session for future use

#### Finding the Session Cookie

After logging in to Collibra in your browser:

**Chrome/Edge:**
1. Press `F12` to open Developer Tools
2. Go to **Application** tab → **Cookies** → select your Collibra domain
3. Find `JSESSIONID` and copy its **Value**

After initial authentication, the cached session is automatically used:
```bash
# Subsequent runs use cached session automatically
./chip --api-url "https://pggm.collibra.com"
```

The session cache is stored at `~/.config/collibra/session_cache.json`.

## VS Code Integration

### 1. Open MCP Settings

Open the Command Palette (`Ctrl+Shift+P`) and search for:
```
MCP: Open User Configuration (JSON)
```

Or navigate directly to the MCP configuration file:
- **Windows:** `%APPDATA%\Code\User\mcp.json` (or `Code - Insiders` for Insiders)

### 2. Add Collibra Server Configuration

Add the following to your `mcp.json` file inside the `"servers"` object:

```json
{
    "servers": {
        "collibra": {
            "type": "stdio",
            "command": "<path-to-chip>/chip.exe",
            "enabled": true,
            "args": ["--api-url", "https://pggm.collibra.com", "--sso-auth"]
        }
    }
}
```

> **Note:** Replace `<path-to-chip>` with the actual path where you built the binary, e.g., `C:\\Projects\\collibra-mcp\\chip.exe`

### 3. Authenticate (First-time only)

Before using the MCP in VS Code, run the chip binary once to cache your SSO session:

```powershell
cd <path-to-chip>
.\chip.exe --api-url "https://pggm.collibra.com" --sso-auth
```

Follow the prompts to complete SSO authentication. After this, VS Code will automatically use the cached session.

### 4. Restart MCP Server

In VS Code, open Command Palette (`Ctrl+Shift+P`) and run:
```
MCP: List Servers
```

Then restart the Collibra server to apply the configuration.

### 5. Verify Installation

To confirm the MCP is working:
1. Open GitHub Copilot Chat in VS Code
2. Type: `@collibra search for assets containing "customer"`
3. You should see results from your Collibra instance

## Troubleshooting

### "Session expired" or authentication errors
Re-run the authentication command:
```bash
./chip --api-url "https://pggm.collibra.com" --sso-auth
```
Then restart the MCP server in VS Code (Command Palette → `MCP: List Servers` → click restart icon next to `collibra`).

### MCP server not showing in VS Code
- Ensure the path in `mcp.json` is correct and uses double backslashes (`\\`) on Windows, e.g., `C:\\path\\to\\chip.exe`
- Check that the chip binary exists at the specified path
- Restart VS Code after updating `mcp.json`

### "JSESSIONID not found" in browser
- Make sure you're fully logged in to Collibra (not on the login page)
- Try refreshing the Collibra page after login
- Clear browser cookies and log in again 
