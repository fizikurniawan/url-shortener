package config

import (
	"os"

	"github.com/gorilla/sessions"
)

var (
	// Store will hold all our session data
	Store *sessions.CookieStore
)

func InitSession() {
	// In production, you should use an environment variable for the secret key
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "your-secret-key-change-this-in-production"
	}

	Store = sessions.NewCookieStore([]byte(sessionKey))

	// Configure session cookie
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,      // For security, prevent JavaScript access
		Secure:   false,     // Set to true if using HTTPS
	}
}
