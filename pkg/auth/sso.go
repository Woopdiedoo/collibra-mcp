package auth

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	// DefaultTimeout is the maximum time to wait for SSO authentication
	DefaultTimeout = 5 * time.Minute
	// SessionCookieName is the name of the Collibra session cookie
	SessionCookieName = "JSESSIONID"
)

// SSOAuthResult contains the result of SSO authentication
type SSOAuthResult struct {
	Cookie    string
	ExpiresAt time.Time
}

// AuthenticateWithSSO opens a browser for the user to authenticate via SSO
// then prompts them to paste the session cookie.
func AuthenticateWithSSO(ctx context.Context, collibraURL string, timeout time.Duration) (*SSOAuthResult, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	slog.Info("Starting SSO authentication...")

	// Open the URL in the default browser (preserves SSO/Intune trust)
	if err := openBrowser(collibraURL); err != nil {
		return nil, fmt.Errorf("failed to open browser: %w", err)
	}

	// Print instructions for the user
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    SSO Authentication                          ║")
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")
	fmt.Println("║  1. Log in to Collibra in the browser that just opened         ║")
	fmt.Println("║  2. Once logged in, press F12 to open Developer Tools          ║")
	fmt.Println("║  3. Go to: Application → Cookies → " + extractDomain(collibraURL))
	fmt.Println("║  4. Find the 'JSESSIONID' cookie and copy its Value            ║")
	fmt.Println("║  5. Paste the cookie value below and press Enter               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Print("Paste JSESSIONID cookie value: ")

	// Read the cookie from stdin
	reader := bufio.NewReader(os.Stdin)
	cookie, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read cookie: %w", err)
	}

	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return nil, fmt.Errorf("no cookie provided")
	}

	slog.Info("SSO authentication successful!")
	return &SSOAuthResult{
		Cookie:    cookie,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}, nil
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// extractDomain extracts the domain from a URL
func extractDomain(url string) string {
	// Remove protocol
	domain := strings.TrimPrefix(url, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	// Remove path
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}

	return domain
}

// ValidateSession checks if a session cookie is still valid by making a test request
func ValidateSession(collibraURL, cookie string) bool {
	// This would make a quick API call to verify the session
	// For now, we just check if the cookie is non-empty
	return cookie != ""
}
