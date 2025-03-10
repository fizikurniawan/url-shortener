package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *sql.DB
	templates *template.Template
	config    *config.Config
}

func NewAuthHandler(db *sql.DB, templates *template.Template, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:        db,
		templates: templates,
		config:    cfg,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := NewPageData("login")
	data.BaseURL = h.config.Env.BaseURL
	data.CurrentTime = h.config.Env.CurrentTime

	if r.Method == http.MethodGet {
		err := h.templates.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Template Error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validasi input
		if username == "" || password == "" {
			data.Error = "Username and password are required"
			h.templates.ExecuteTemplate(w, "layout.html", data)
			return
		}

		var user model.User
		var hashedPassword string

		err := h.db.QueryRow(
			"SELECT id, username, password, is_active FROM users WHERE username = $1",
			username,
		).Scan(&user.ID, &user.Username, &hashedPassword, &user.IsActive)

		if err != nil {
			log.Printf("Database error: %v", err)
			data.Error = "Invalid username or password"
			data.Data["Username"] = username
			h.templates.ExecuteTemplate(w, "layout.html", data)
			return
		}

		// Compare password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Printf("Password comparison error: %v", err)
			data.Error = "Invalid username or password"
			data.Data["Username"] = username
			h.templates.ExecuteTemplate(w, "layout.html", data)
			return
		}

		// if user inactivated
		if !user.IsActive {
			data.Error = "User is inactive"
			data.Data["Username"] = username
			h.templates.ExecuteTemplate(w, "layout.html", data)
			return
		}

		// Get session
		session, err := h.config.Store.Get(r, "user-session")
		if err != nil {
			log.Printf("Session error: %v", err)
			data.Error = "Session Error"
			h.templates.ExecuteTemplate(w, "layout.html", data)
			return
		}

		// Set session values with explicit typing
		session.Values["user_id"] = int64(user.ID)
		session.Values["username"] = string(user.Username)

		// Save session
		err = session.Save(r, w)
		if err != nil {
			log.Printf("Session save error: %v", err)
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		// Redirect after successful login
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := h.config.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Session Error", http.StatusInternalServerError)
		return
	}

	// Clear session
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
