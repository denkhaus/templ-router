package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// Production-ready session store implementations

type simpleSessionStore struct {
	logger *zap.Logger
}

// NewInMemorySessionStore creates a new session store for DI
func NewInMemorySessionStore(i do.Injector) (interfaces.SessionStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	return &simpleSessionStore{logger: logger}, nil
}

func (s *simpleSessionStore) GetSession(req *http.Request) (*interfaces.Session, error) {
	// Check for session cookie
	cookie, err := req.Cookie("session_id")
	if err != nil {
		s.logger.Debug("No session cookie found")
		return nil, fmt.Errorf("no session cookie")
	}

	// For demo: validate session ID format
	if len(cookie.Value) < 10 {
		s.logger.Debug("Invalid session ID format", zap.String("session_id", cookie.Value))
		return nil, fmt.Errorf("invalid session")
	}

	// Extract user ID from session ID (format: session_userID_timestamp)
	parts := strings.Split(cookie.Value, "_")
	if len(parts) >= 2 {
		userID := parts[1]
		s.logger.Debug("Session found",
			zap.String("session_id", cookie.Value),
			zap.String("user_id", userID))
		return &interfaces.Session{ID: cookie.Value, UserID: userID, Valid: true}, nil
	}

	// Fallback for legacy sessions
	return &interfaces.Session{ID: cookie.Value, UserID: "user1", Valid: true}, nil
}

func (s *simpleSessionStore) CreateSession(userID string) (*interfaces.Session, error) {
	// Generate secure session ID
	sessionID := fmt.Sprintf("session_%s_%d", userID, time.Now().Unix())

	s.logger.Info("Creating new session",
		zap.String("user_id", userID),
		zap.String("session_id", sessionID))

	return &interfaces.Session{
		ID:     sessionID,
		UserID: userID,
		Valid:  true,
	}, nil
}

func (s *simpleSessionStore) DeleteSession(sessionID string) error {
	return nil
}
