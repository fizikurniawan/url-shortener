package middleware

import (
	"context"
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/model"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthRequired(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := cfg.Store.Get(r, "user-session")
		if err != nil {
			log.Printf("Session error in middleware: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get user data from session with type assertions
		userIDInterface := session.Values["user_id"]
		usernameInterface := session.Values["username"]

		userID, ok1 := userIDInterface.(int64)
		username, ok2 := usernameInterface.(string)

		if !ok1 || !ok2 {
			log.Printf("Invalid session data: userID(%T) = %v, username(%T) = %v",
				userIDInterface, userIDInterface,
				usernameInterface, usernameInterface)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Create user object
		user := &model.User{
			ID:       userID,
			Username: username,
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func OptionalAuth(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := cfg.Store.Get(r, "user-session")
		if err != nil {
			log.Printf("Session error in middleware: %v", err)
		}

		userIDInterface := session.Values["user_id"]
		usernameInterface := session.Values["username"]

		userID, ok1 := userIDInterface.(int64)
		username, ok2 := usernameInterface.(string)

		if ok1 && ok2 {
			user := &model.User{
				ID:       userID,
				Username: username,
			}
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	}
}
