package handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/service"
	"url-shortener/internal/util"
)

type URLHandler struct {
	urlService *service.URLService
	templates  *template.Template
	config     *config.Config
}

func NewURLHandler(urlService *service.URLService, templates *template.Template, cfg *config.Config) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		templates:  templates,
		config:     cfg,
	}
}

func (h *URLHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	user := util.GetUserFromContext(r.Context())
	pageData := NewPageData("index") // Specify which template to use
	pageData.BaseURL = h.config.Env.BaseURL
	pageData.CurrentTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	pageData.User = user
	pageData.Data["URL"] = nil

	if r.Method == http.MethodGet {
		err := h.templates.ExecuteTemplate(w, "layout.html", pageData)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		longURL := r.FormValue("url")

		if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
			longURL = "https://" + longURL
		}

		var userID *int64
		var customCode string

		if user != nil {
			userID = &user.ID
			customCode = r.FormValue("custom_code")
		}

		url, err := h.urlService.CreateShortURL(longURL, userID, customCode)
		if err != nil {
			pageData.Error = err.Error()
		} else {
			pageData.Data["URL"] = url
		}

		err = h.templates.ExecuteTemplate(w, "layout.html", pageData)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	}
}

func (h *URLHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	// Add cache control headers
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	shortCode := strings.Trim(r.URL.Path, "/")

	if shortCode == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	url, err := h.urlService.GetURL(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusMovedPermanently)
}

func (h *URLHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	user := util.GetUserFromContext(r.Context())
	pageData := NewPageData("dashboard")
	pageData.BaseURL = h.config.Env.BaseURL
	pageData.CurrentTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	pageData.User = user

	// Get user from session
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user's URLs
	urls, err := h.urlService.GetUserURLs(user.ID)
	if err != nil {
		log.Printf("Error getting user URLs: %v", err)
		pageData.Error = "Failed to load URLs"
		h.templates.ExecuteTemplate(w, "layout", pageData)
		return
	}

	pageData.Data["URLs"] = urls
	err = h.templates.ExecuteTemplate(w, "layout.html", pageData)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
