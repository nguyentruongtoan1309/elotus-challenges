package utils

import (
	"sync"
	"time"
)

// TokenBlacklist manages revoked tokens
type TokenBlacklist struct {
	tokens map[string]time.Time
	mutex  sync.RWMutex
}

// NewTokenBlacklist creates a new TokenBlacklist instance
func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
}

// RevokeToken adds a token to the blacklist
func (tb *TokenBlacklist) RevokeToken(tokenString string) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.tokens[tokenString] = time.Now()
	tb.cleanup()
}

// IsRevoked checks if a token has been revoked
func (tb *TokenBlacklist) IsRevoked(tokenString string) bool {
	tb.mutex.RLock()
	defer tb.mutex.RUnlock()
	_, exists := tb.tokens[tokenString]
	return exists
}

// cleanup removes expired tokens from the blacklist
func (tb *TokenBlacklist) cleanup() {
	now := time.Now()
	for token, revokedAt := range tb.tokens {
		// Remove tokens that were revoked more than 24 hours ago (token expiry time)
		if now.Sub(revokedAt) > 24*time.Hour {
			delete(tb.tokens, token)
		}
	}
}

// Global token blacklist instance
var globalTokenBlacklist = NewTokenBlacklist()

// GetTokenBlacklist returns the global token blacklist instance
func GetTokenBlacklist() *TokenBlacklist {
	return globalTokenBlacklist
}
