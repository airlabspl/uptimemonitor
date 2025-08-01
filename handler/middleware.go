package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"uptimemonitor"
)

type contextKey string

const (
	sessionContextKey contextKey = "session"
	userContextKey    contextKey = "user"
)

func (m *Handler) Installed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, err := m.Store.CountUsers(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if count == 0 && r.Method == http.MethodGet {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}

		if count == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Handler) UserFromCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil || c.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		session, err := m.Store.GetSessionByUuid(r.Context(), c.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey, session)
		ctx = context.WithValue(ctx, userContextKey, session.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Handler) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Context().Value(sessionContextKey)
		if value == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, ok := value.(uptimemonitor.Session)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Handler) Guest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Context().Value(sessionContextKey)
		if value != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Handler) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Handler) NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/static") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}
