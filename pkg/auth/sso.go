package auth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
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

// AuthenticateWithSSO opens a browser window for the user to authenticate via SSO
// and captures the session cookie after successful authentication.
func AuthenticateWithSSO(ctx context.Context, collibraURL string, timeout time.Duration) (*SSOAuthResult, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	slog.Info("Starting SSO authentication...")
	slog.Info(fmt.Sprintf("Opening browser to: %s", collibraURL))

	// Find or download a browser
	path, _ := launcher.LookPath()
	if path == "" {
		slog.Info("No browser found, downloading Chromium...")
	}

	// Launch browser in headed mode so user can interact with SSO
	l := launcher.New().
		Headless(false).
		Set("disable-gpu").
		Set("no-first-run").
		Set("no-default-browser-check")

	controlURL, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(controlURL)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Navigate to Collibra (will trigger SSO redirect)
	page, err := browser.Page(proto.TargetCreateTarget{URL: collibraURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	slog.Info("Waiting for SSO authentication to complete...")
	slog.Info("Please log in using the browser window that opened.")

	// Wait for the session cookie to appear (indicates successful auth)
	result, err := waitForSessionCookie(ctx, page, collibraURL)
	if err != nil {
		return nil, err
	}

	slog.Info("SSO authentication successful!")
	return result, nil
}

// waitForSessionCookie polls for the session cookie until it appears or timeout
func waitForSessionCookie(ctx context.Context, page *rod.Page, collibraURL string) (*SSOAuthResult, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Extract the domain from the URL for cookie matching
	domain := extractDomain(collibraURL)

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("SSO authentication timed out - please try again")
		case <-ticker.C:
			// Check if we're back on the Collibra domain
			info, err := page.Info()
			if err != nil {
				continue
			}

			// Only check cookies when we're on the Collibra domain
			if !strings.Contains(info.URL, domain) {
				continue
			}

			// Get all cookies for the page
			cookies, err := page.Cookies([]string{collibraURL})
			if err != nil {
				continue
			}

			// Look for the session cookie
			for _, cookie := range cookies {
				if cookie.Name == SessionCookieName {
					expiresAt := time.Now().Add(30 * time.Minute) // Default session timeout
					if cookie.Expires > 0 {
						expiresAt = time.Unix(int64(cookie.Expires), 0)
					}

					return &SSOAuthResult{
						Cookie:    cookie.Value,
						ExpiresAt: expiresAt,
					}, nil
				}
			}
		}
	}
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
