package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/denkhaus/templ-router/demo/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// demoSessionStore provides a default in-memory session store implementation for the demo
// This shows how to implement a session store for production applications
type demoSessionStore struct {
	logger        *zap.Logger
	sessions      map[string]*interfaces.Session
	mutex         sync.RWMutex
	sessionExpiry time.Duration
	cookieName    string
}

// NewDemoSessionStore creates a new demo session store for DI
func NewDemoSessionStore(i do.Injector) (interfaces.SessionStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)

	store := &demoSessionStore{
		logger:        logger,
		sessions:      make(map[string]*interfaces.Session),
		mutex:         sync.RWMutex{},
		cookieName:    "session_id",   // Hardcoded for demo
		sessionExpiry: 24 * time.Hour, // Hardcoded for demo
	}

	// Start cleanup routine for expired sessions
	go store.cleanupExpiredSessions()

	return store, nil
}

// GetSession retrieves a session from the request
func (s *demoSessionStore) GetSession(req *http.Request) (*interfaces.Session, error) {
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
func (s *demoSessionStore) CreateSession(userID string) (*interfaces.Session, error) {
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
		zap.String("session_id", sessionID),
		zap.String("user_id", userID))

	return session, nil
}

// DeleteSession deletes a session
func (s *demoSessionStore) DeleteSession(sessionID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, sessionID)
	s.logger.Info("Session deleted", zap.String("session_id", sessionID))
	return nil
}

// GetSessionByID retrieves a session by its ID (direct access)
func (s *demoSessionStore) GetSessionByID(sessionID string) (*interfaces.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

// ExtendSession extends the expiry time of an existing session
func (s *demoSessionStore) ExtendSession(sessionID string, duration time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	// Check if session is already expired
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, sessionID)
		return fmt.Errorf("session expired and cannot be extended")
	}

	// Extend the session expiry
	session.ExpiresAt = time.Now().Add(duration)
	s.sessions[sessionID] = session

	s.logger.Info("Session extended",
		zap.String("session_id", sessionID),
		zap.String("user_id", session.UserID),
		zap.Duration("duration", duration),
		zap.Time("new_expiry", session.ExpiresAt))

	return nil
}

// generateSessionID generates a cryptographically secure session ID
func (s *demoSessionStore) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// cleanupExpiredSessions removes expired sessions from memory
// This runs in a background goroutine
func (s *demoSessionStore) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
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
				zap.Int("expired_count", len(expiredSessions)),
				zap.Int("remaining_sessions", len(s.sessions)))
		}
	}
}
