package tools

import (
	"context"
	"fmt"
	"net/http"
	"os"

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

		// Get the path to the currently running chip executable
		execPath, err := os.Executable()
		if err != nil {
			execPath = "chip" // Fallback to assuming it's in PATH
		}

		instructions := fmt.Sprintf(`AUTHENTICATION REQUIRED - AI AGENT INSTRUCTIONS
================================================

You (the AI agent) must perform the following steps to authenticate with Collibra.
Do NOT just tell the user these instructions - YOU must execute them.

STEP 1: RUN THE AUTHENTICATION COMMAND
--------------------------------------
Use your run_in_terminal tool to execute this command as a BACKGROUND process (isBackground=true):

    %s --api-url "%s" --sso-auth

IMPORTANT: You MUST set isBackground=true so you don't get stuck waiting for the terminal.

STEP 2: TELL THE USER WHAT TO DO
--------------------------------
After running the command, immediately tell the user:

"I've opened your browser for Collibra authentication. Please:
1. Log in to Collibra in the browser window that just opened
2. Once logged in, press F12 to open Developer Tools
3. Go to: Application → Cookies → %s
4. Find the 'JSESSIONID' cookie and copy its Value
5. Click on the terminal in VS Code and paste the cookie value, then press Enter

Let me know once you've completed this and I'll retry the search!"

STEP 3: WAIT FOR USER CONFIRMATION
----------------------------------
Wait for the user to confirm they have completed authentication before retrying any Collibra operations.
The session will be cached automatically for future use.

DO NOT just display these instructions to the user - YOU must run the terminal command yourself!
`, execPath, collibraHost, collibraHost)

		return AuthHelpOutput{
			Instructions: instructions,
			Status:       "Authentication help requested - please follow the instructions above",
		}, nil
	}
}
