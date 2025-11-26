package services

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/denkhaus/templ-router/demo/pkg/interfaces"
	"go.uber.org/zap"
)

// CustomSessionStore demonstrates a simple in-memory session store implementation
// This shows how to create a custom SessionStore for your application
type CustomSessionStore struct {
	logger        *zap.Logger
	sessionExpiry time.Duration
	cookieName    string
	sessions      map[string]*interfaces.Session
	mutex         sync.RWMutex
}

// NewCustomSessionStore creates a new custom in-memory session store
func NewCustomSessionStore(configService any, logger *zap.Logger) (interfaces.SessionStore, error) {
	logger.Info("Custom in-memory session store initialized")

	return &CustomSessionStore{
		logger:        logger,
		sessionExpiry: 24 * time.Hour, // Hardcoded for demo
		cookieName:    "session_id",   // Hardcoded for demo
		sessions:      make(map[string]*interfaces.Session),
		mutex:         sync.RWMutex{},
	}, nil
}

// GetSession retrieves a session from memory
func (s *CustomSessionStore) GetSession(req *http.Request) (*interfaces.Session, error) {
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

// CreateSession creates a new session in memory
func (s *CustomSessionStore) CreateSession(userID string) (*interfaces.Session, error) {
	sessionID := fmt.Sprintf("custom_sess_%d_%d", time.Now().UnixNano(), len(s.sessions)+1)

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

	s.logger.Info("Session created in custom store",
		zap.String("session_id", sessionID),
		zap.String("user_id", userID),
		zap.Duration("expires_in", s.sessionExpiry))

	return session, nil
}

// DeleteSession deletes a session from memory
func (s *CustomSessionStore) DeleteSession(sessionID string) error {
	s.mutex.Lock()
	delete(s.sessions, sessionID)
	s.mutex.Unlock()

	s.logger.Info("Session deleted from custom store", zap.String("session_id", sessionID))
	return nil
}

// GetSessionByID retrieves a session by its ID (direct access)
func (s *CustomSessionStore) GetSessionByID(sessionID string) (*interfaces.Session, error) {
	s.mutex.RLock()
	session, exists := s.sessions[sessionID]
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

// ExtendSession extends the expiry time of an existing session
func (s *CustomSessionStore) ExtendSession(sessionID string, duration time.Duration) error {
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

	s.logger.Info("Session extended in custom store",
		zap.String("session_id", sessionID),
		zap.String("user_id", session.UserID),
		zap.Duration("duration", duration),
		zap.Time("new_expiry", session.ExpiresAt))

	return nil
}

// GetSessionCount returns the number of active sessions (for monitoring)
func (s *CustomSessionStore) GetSessionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.sessions)
}
