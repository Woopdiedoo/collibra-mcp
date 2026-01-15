package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/collibra/chip/pkg/chip"
)

type AuthHelpInput struct {
	// No input required - this tool just provides authentication instructions
}

type AuthHelpOutput struct {
	Instructions string `json:"instructions" jsonschema:"Step-by-step instructions for authenticating with the Collibra MCP server"`
	Status       string `json:"status" jsonschema:"Current authentication status or the error that triggered this help request"`
}

func NewAuthHelpTool(_ *http.Client) *chip.Tool[AuthHelpInput, AuthHelpOutput] {
	return &chip.Tool[AuthHelpInput, AuthHelpOutput]{
		Name: "auth_help",
		Description: `Get authentication help and instructions for the Collibra MCP server.

WHEN TO USE THIS TOOL:
- When you receive HTTP 401 (Unauthorized) or 403 (Forbidden) errors from other Collibra tools
- When the user asks about authentication or login issues
- When any Collibra tool returns an authentication-related error

This tool will provide step-by-step instructions that you can share with the user to resolve authentication issues.`,
		Handler: handleAuthHelp(),
	}
}

func handleAuthHelp() chip.ToolHandlerFunc[AuthHelpInput, AuthHelpOutput] {
	return func(ctx context.Context, input AuthHelpInput) (AuthHelpOutput, error) {
		collibraHost, _ := chip.GetCollibraHost(ctx)
		if collibraHost == "" {
			collibraHost = "https://pggm.collibra.com"
		}

		instructions := fmt.Sprintf(`AUTHENTICATION INSTRUCTIONS FOR COLLIBRA MCP SERVER
================================================

The Collibra MCP server uses browser-based SSO authentication. Here's how to authenticate:

STEP 1: RESTART THE MCP SERVER WITH SSO FLAG
--------------------------------------------
The user needs to restart the Collibra MCP server with the --sso flag enabled.
This can be done by:
- Adding "ssoAuth": true to their chip configuration file, OR
- Running the server with: chip --sso

STEP 2: FOLLOW THE BROWSER AUTHENTICATION FLOW
----------------------------------------------
When the server starts with SSO enabled:
1. A browser window will automatically open to: %s
2. The user should log in using their corporate SSO credentials (e.g., Okta, Azure AD, etc.)
3. Once logged in, the user needs to:
   a. Press F12 to open Developer Tools
   b. Go to: Application → Cookies → %s
   c. Find the 'JSESSIONID' cookie and copy its Value
   d. Paste the cookie value in the terminal where the MCP server is running

STEP 3: SESSION CACHING
-----------------------
After successful authentication:
- The session is cached locally for future use
- Subsequent MCP server starts will use the cached session automatically
- Sessions typically expire after 30 minutes of inactivity

ALTERNATIVE: DIRECT COOKIE CONFIGURATION
----------------------------------------
If the user already has a valid JSESSIONID cookie, they can configure it directly:
- Set CHIP_API_COOKIE environment variable: CHIP_API_COOKIE="JSESSIONID=<cookie_value>"
- Or add to configuration file: "cookie": "JSESSIONID=<cookie_value>"

TROUBLESHOOTING
---------------
- If authentication keeps failing, the cached session may have expired
- Delete the session cache file and restart with --sso to re-authenticate
- Session cache is typically stored at: ~/.chip/session_cache.json

TELL THE USER:
"To authenticate with Collibra, please restart the MCP server with SSO enabled (use --sso flag or set ssoAuth: true in config). A browser will open for you to log in, then follow the prompts to complete authentication."
`, collibraHost, collibraHost)

		return AuthHelpOutput{
			Instructions: instructions,
			Status:       "Authentication help requested - please follow the instructions above",
		}, nil
	}
}
