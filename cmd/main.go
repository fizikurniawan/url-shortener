package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()
	defer cfg.Close()

	// Log startup info
	log.Printf("Server starting with configuration:")
	log.Printf("Current Time (UTC): %s", cfg.Env.CurrentTime)
	log.Printf("Server Port: %s", cfg.Env.ServerPort)
	log.Printf("Base URL: %s", cfg.Env.BaseURL)

	// Initialize repositories and services
	urlRepo := repository.NewURLRepository(cfg.DB)
	urlService := service.NewURLService(urlRepo)

	// Initialize templates with functions
	templates := template.Must(template.New("").ParseFiles(
		"templates/layout.html",
		"templates/index.html",
		"templates/login.html",
		"templates/register.html",
		"templates/dashboard.html",
	))

	// Initialize handlers
	urlHandler := handler.NewURLHandler(urlService, templates, cfg)
	authHandler := handler.NewAuthHandler(cfg.DB, templates, cfg)
	userHandler := handler.NewUserHandler(cfg.DB, templates, cfg)

	// Create new mux
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Auth routes
	mux.HandleFunc("/login", middleware.OptionalAuth(cfg, authHandler.HandleLogin))
	mux.HandleFunc("/logout", authHandler.HandleLogout)
	mux.HandleFunc("/register", userHandler.HandleRegister)

	// Protected routes
	mux.HandleFunc("/dashboard", middleware.AuthRequired(cfg, urlHandler.HandleDashboard))

	// URL routes
	mux.HandleFunc("/", middleware.OptionalAuth(cfg, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			urlHandler.HandleIndex(w, r)
			return
		}
		urlHandler.HandleRedirect(w, r)
	}))

	// Create server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Env.ServerPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server
	log.Printf("Server listening on port %s", cfg.Env.ServerPort)
	log.Fatal(server.ListenAndServe())
}
