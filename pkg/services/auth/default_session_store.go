package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// inMemmorySessionStoreImpl provides a default in-memory session store implementation
// Users can replace this with Redis, database-backed, or other implementations
type inMemmorySessionStoreImpl struct {
	logger        *zap.Logger
	sessions      map[string]*interfaces.Session
	mutex         sync.RWMutex
	configService interfaces.ConfigService
	sessionExpiry time.Duration
	cookieName    string
}

// NewInMemorySessionStore creates a new default session store for DI
func NewInMemorySessionStore(i do.Injector) (interfaces.SessionStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)

	store := &inMemmorySessionStoreImpl{
		logger:        logger,
		sessions:      make(map[string]*interfaces.Session),
		mutex:         sync.RWMutex{},
		configService: configService,
		cookieName:    configService.GetSessionCookieName(),
		sessionExpiry: configService.GetSessionExpiry(),
	}

	// Start cleanup routine for expired sessions
	go store.cleanupExpiredSessions()

	return store, nil
}

// GetSession retrieves a session from the request
func (s *inMemmorySessionStoreImpl) GetSession(req *http.Request) (*interfaces.Session, error) {
	// Get session ID from cookie
	cookie, err := req.Cookie(s.cookieName)
	if err != nil {
		return nil, fmt.Errorf("no session cookie found")
	}

	s.mutex.RLock()
	session, exists := s.sessions[cookie.Value]
	s.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(session.ID)
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

// CreateSession creates a new session for a user
func (s *inMemmorySessionStoreImpl) CreateSession(userID string) (*interfaces.Session, error) {
	sessionID, err := s.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &interfaces.Session{
		ID:        sessionID,
		UserID:    userID,
		Valid:     true,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionExpiry),
	}

	s.mutex.Lock()
	s.sessions[sessionID] = session
	s.mutex.Unlock()

	s.logger.Info("Session created",
		zap.String(s.cookieName, sessionID),
		zap.String("user_id", userID))

	return session, nil
}

// DeleteSession deletes a session
func (s *inMemmorySessionStoreImpl) DeleteSession(sessionID string) error {
	s.mutex.Lock()
	delete(s.sessions, sessionID)
	s.mutex.Unlock()

	s.logger.Info("Session deleted", zap.String(s.cookieName, sessionID))
	return nil
}

// generateSessionID generates a cryptographically secure session ID
func (s *inMemmorySessionStoreImpl) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// cleanupExpiredSessions runs a background routine to clean up expired sessions
func (s *inMemmorySessionStoreImpl) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		expiredSessions := []string{}

		for sessionID, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				expiredSessions = append(expiredSessions, sessionID)
			}
		}

		for _, sessionID := range expiredSessions {
			delete(s.sessions, sessionID)
		}
		s.mutex.Unlock()

		if len(expiredSessions) > 0 {
			s.logger.Info("Cleaned up expired sessions",
				zap.Int("count", len(expiredSessions)))
		}
	}
}
