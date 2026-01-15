# Configuration Guide

The Collibra MCP Server uses cookie-based SSO authentication.

## Authentication

### First-time Setup

Run the server with the `--sso-auth` flag to authenticate:

```bash
./chip --api-url "https://pggm.collibra.com" --sso-auth
```

This will:
1. Open your default browser to the Collibra login page
2. Allow you to complete SSO authentication
3. Prompt you to paste the session cookie from browser DevTools
4. Cache the session for future use

### Subsequent Runs

After initial authentication, the cached session is automatically used:

```bash
./chip --api-url "https://pggm.collibra.com" --sso-auth
```

## Command Line Options

```
--api-url           Collibra API URL (required)
--sso-auth          Enable browser-based SSO authentication
--sso-timeout       Timeout in seconds for SSO authentication (default: 300)
--sso-cache-path    Custom path to cache SSO session
--mode              Server mode: 'stdio' (default) or 'http'
--port              HTTP server port (default: 8080, only used in http mode)
--skip-tls-verify   Skip TLS certificate verification (for development only)
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `COLLIBRA_MCP_API_URL` | Collibra API base URL |
| `COLLIBRA_MCP_SSO_AUTH` | Enable SSO authentication (true/false) |
| `COLLIBRA_MCP_SSO_TIMEOUT` | SSO authentication timeout in seconds |
| `COLLIBRA_MCP_SSO_CACHE_PATH` | Path to cache SSO session |
| `COLLIBRA_MCP_MODE` | Server mode: stdio or http |
| `COLLIBRA_MCP_HTTP_PORT` | HTTP server port |

## Session Cache

The SSO session is cached at:
- `~/.config/collibra/session_cache.json`

The cache includes the session cookie and expiration time. When the session expires, re-run with `--sso-auth` to re-authenticate.
