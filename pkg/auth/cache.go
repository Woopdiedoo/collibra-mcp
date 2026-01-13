package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// CachedSession represents a cached SSO session
type CachedSession struct {
	Cookie    string    `json:"cookie"`
	ExpiresAt time.Time `json:"expires_at"`
	URL       string    `json:"url"`
}

// SessionCache handles caching of SSO session cookies
type SessionCache struct {
	cachePath string
}

// NewSessionCache creates a new session cache
func NewSessionCache(cachePath string) *SessionCache {
	if cachePath == "" {
		// Default to user's config directory
		configDir, err := os.UserConfigDir()
		if err != nil {
			configDir = os.TempDir()
		}
		cachePath = filepath.Join(configDir, "collibra", "session_cache.json")
	}
	return &SessionCache{cachePath: cachePath}
}

// Load retrieves a cached session if it exists and is still valid
func (c *SessionCache) Load(collibraURL string) (*CachedSession, error) {
	data, err := os.ReadFile(c.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cache exists
		}
		return nil, fmt.Errorf("failed to read session cache: %w", err)
	}

	var session CachedSession
	if err := json.Unmarshal(data, &session); err != nil {
		slog.Warn("Failed to parse session cache, will re-authenticate")
		return nil, nil
	}

	// Check if the session is for the same URL
	if session.URL != collibraURL {
		slog.Info("Cached session is for a different URL, will re-authenticate")
		return nil, nil
	}

	// Check if the session has expired (with 5 minute buffer)
	if time.Now().Add(5 * time.Minute).After(session.ExpiresAt) {
		slog.Info("Cached session has expired, will re-authenticate")
		return nil, nil
	}

	slog.Info("Using cached SSO session")
	return &session, nil
}

// Save stores a session to the cache
func (c *SessionCache) Save(result *SSOAuthResult, collibraURL string) error {
	session := CachedSession{
		Cookie:    result.Cookie,
		ExpiresAt: result.ExpiresAt,
		URL:       collibraURL,
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(c.cachePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Write with restricted permissions (owner only)
	if err := os.WriteFile(c.cachePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write session cache: %w", err)
	}

	slog.Info(fmt.Sprintf("Session cached to: %s", c.cachePath))
	return nil
}

// Clear removes the cached session
func (c *SessionCache) Clear() error {
	if err := os.Remove(c.cachePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear session cache: %w", err)
	}
	return nil
}

// GetCachePath returns the path to the cache file
func (c *SessionCache) GetCachePath() string {
	return c.cachePath
}
