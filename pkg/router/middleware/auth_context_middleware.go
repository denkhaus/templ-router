package middleware

import (
	"context"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// AuthContextMiddleware adds user information to the request context
type AuthContextMiddleware struct {
	sessionStore interfaces.SessionStore
	userStore    interfaces.UserStore
	logger       *zap.Logger
}

// NewAuthContextMiddleware creates a new auth context middleware
func NewAuthContextMiddleware(i do.Injector) (*AuthContextMiddleware, error) {
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	userStore := do.MustInvoke[interfaces.UserStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &AuthContextMiddleware{
		sessionStore: sessionStore,
		userStore:    userStore,
		logger:       logger,
	}, nil
}

// Middleware returns the HTTP middleware function
func (acm *AuthContextMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		acm.logger.Info("Auth middleware processing request", zap.String("path", r.URL.Path))

		// Try to get session from request
		session, err := acm.sessionStore.GetSession(r)
		if err != nil {
			// No session found, continue without user
			acm.logger.Info("No session found", zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		acm.logger.Info("Session found",
			zap.String("session_id", session.ID),
			zap.String("user_id", session.UserID),
		)

		// Get user from session
		user, err := acm.userStore.GetUserByID(session.UserID)
		if err != nil {
			acm.logger.Warn("Failed to get user from session",
				zap.String("session_id", session.ID),
				zap.String("user_id", session.UserID),
				zap.Error(err))
			next.ServeHTTP(w, r)
			return
		}

		// Add user to context using the correct key constant
		ctx = context.WithValue(ctx, shared.UserContextKey, user)
		r = r.WithContext(ctx)

		acm.logger.Debug("User added to context",
			zap.String("user_id", user.GetID()),
			zap.String("email", user.GetEmail()))

		next.ServeHTTP(w, r)
	})
}
