package auth

import (
	"context"
	"net/http"
	"strings"

	"portal/internal/auth/store"
)

const SessionCookie = "portal_session"

func Auth(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions ||
				strings.HasPrefix(r.URL.Path, "/proxy/") ||
				strings.Contains(r.URL.Path, "/_next/") ||
				strings.Contains(r.URL.Path, "/static/") ||
				strings.HasPrefix(r.URL.Path, "/login") ||
				strings.HasPrefix(r.URL.Path, "/csrf") ||
				r.URL.Path == "/api/health" {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(SessionCookie)
			if err != nil || cookie.Value == "" {
				redirectLogin(w, r, "missing")
				return
			}

			username, valid := s.ValidateSession(cookie.Value)
			if !valid {
				redirectLogin(w, r, "expired")
				return
			}

			ctx := WithUser(r.Context(), username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func redirectLogin(w http.ResponseWriter, r *http.Request, reason string) {
	url := "/login"
	if reason != "" && reason != "missing" {
		url += "?reason=" + reason
	}
	http.Redirect(w, r, url, http.StatusFound)
}

type contextKey string

const userKey contextKey = "portal_user"

func WithUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userKey, username)
}

func GetUser(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userKey).(string)
	return username, ok
}
