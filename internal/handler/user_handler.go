package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"
	"url-shortener/internal/config"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db        *sql.DB
	templates *template.Template
	config    *config.Config
}

func NewUserHandler(db *sql.DB, templates *template.Template, cfg *config.Config) *UserHandler {
	return &UserHandler{
		db:        db,
		templates: templates,
		config:    cfg,
	}
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	pageData := NewPageData("register")
	pageData.BaseURL = h.config.Env.BaseURL
	pageData.CurrentTime = time.Now().UTC().Format("2006-01-02 15:04:05")

	if r.Method == http.MethodGet {
		err := h.templates.ExecuteTemplate(w, "layout.html", pageData)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Validasi input
		if username == "" || password == "" {
			pageData.Error = "Username and password are required"
			h.templates.ExecuteTemplate(w, "layout.html", pageData)
			return
		}

		if password != confirmPassword {
			pageData.Error = "Passwords do not match"
			pageData.Data["Username"] = username
			h.templates.ExecuteTemplate(w, "layout.html", pageData)
			return
		}

		// Check if username already exists
		var exists bool
		err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if exists {
			pageData.Error = "Username already exists"
			pageData.Data["Username"] = username
			h.templates.ExecuteTemplate(w, "layout.html", pageData)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Insert new user
		_, err = h.db.Exec(
			"INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3)",
			username,
			hashedPassword,
			time.Now().UTC(),
		)

		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}
