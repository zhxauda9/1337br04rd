package middleware

import (
	"context"
	"fmt"
	"net/http"

	"1337b04rd/internal/domain"
)

type contextKey string

const sessionKey contextKey = "session"

// InjectSessionMiddleware injects the session into the request context
func InjectSessionMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve session ID from cookie
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized: session cookie missing", http.StatusUnauthorized)
				return
			}

			// Retrieve session name from the session_name cookie
			nameCookie, err := r.Cookie("session_name")
			if err != nil || nameCookie.Value == "" {
				http.Error(w, "Unauthorized: session name cookie missing", http.StatusUnauthorized)
				return
			}

			// Create a session object with both ID and name
			session := &domain.Session{
				ID:   cookie.Value,
				Name: nameCookie.Value,
			}

			// Store session in context
			ctx := context.WithValue(r.Context(), sessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSession(ctx context.Context) (*domain.Session, error) {
	session, ok := ctx.Value(sessionKey).(*domain.Session)
	if !ok || session == nil {
		return nil, fmt.Errorf("no session found")
	}
	return session, nil
}
