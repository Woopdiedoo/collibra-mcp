# Collibra MCP Server

A Model Context Protocol (MCP) server that provides AI agents with access to Collibra Data Governance Center capabilities including data asset discovery, business glossary queries, and detailed asset information retrieval.

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

- Access to a Collibra Data Governance Center instance
- Valid Collibra credentials

### Installation

#### Option A: Download Prebuilt Binary (Recommended)

1. **Download the latest release:**
   - Go to the [GitHub Releases page](../../releases)
   - Download the appropriate binary for your platform:
     - `chip-linux-amd64` - Linux (Intel/AMD 64-bit)
     - `chip-linux-arm64` - Linux (ARM 64-bit)
     - `chip-mac-amd64` - macOS (Intel)
     - `chip-mac-arm64` - macOS (Apple Silicon)
     - `chip-windows-amd64.exe` - Windows (Intel/AMD 64-bit)
     - `chip-windows-arm64.exe` - Windows (ARM 64-bit)

3. **Optional: Move to your PATH:**
   ```bash
   # Linux/macOS
   sudo mv chip-* /usr/local/bin/mcp-server
   
   # Or add to your user bin directory
   mv chip-* ~/.local/bin/mcp-server
   ```

#### Option B: Build from Source
   ```bash
   git clone <repository-url>
   cd chip
   go mod download
   go build -o .build/chip ./cmd/chip

   # Run the build binary
   ./.build/chip
   ```

## Running and Configuration

### Authentication

The server uses cookie-based authentication with SSO support:

```bash
# First-time setup: authenticate via browser
./chip --api-url "https://your-collibra-instance.com" --sso-auth
```

This will:
1. Open your default browser to the Collibra login page
2. Allow you to complete SSO authentication
3. Prompt you to paste the session cookie from browser DevTools
4. Cache the session for future use

After initial authentication, the cached session is automatically used:
```bash
# Subsequent runs use cached session automatically
./chip --api-url "https://your-collibra-instance.com"
```

The session cache is stored at `~/.config/collibra/session_cache.json`.

**For detailed configuration instructions, see [CONFIG.md](docs/CONFIG.md).**

## Integration with MCP Clients

This server is compatible with any MCP client. Refer to your MCP client's documentation for server configuration. 

### VS Code / VS Code Insiders

```json
// User settings: mcp.json or .vscode/mcp.json
{
    "servers": {
        "collibra": {
            "type": "stdio",
            "command": "/path/to/chip",
            "args": ["--api-url", "https://your-collibra-instance.com"]
        }
    }
}
```

**Note:** For SSO authentication, run the chip binary once manually with `--sso-auth` to cache your session. VS Code will then automatically use the cached session. 
