package config

import (
	"database/sql"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite" // Import SQLite driver
)

func init() {
	// Register types for session serialization
	gob.Register(int64(0))
	gob.Register("")
	gob.Register(map[string]interface{}{})
}

type Environment struct {
	ServerPort  string
	BaseURL     string
	CurrentUser string
	CurrentTime string
}

type Config struct {
	DB    *sql.DB
	Store sessions.Store
	Env   *Environment
}

func NewConfig() *Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Setup environment
	env := &Environment{
		ServerPort:  getEnvOrDefault("PORT", "8080"),
		BaseURL:     getEnvOrDefault("BASE_URL", "http://localhost:8080"),
		CurrentUser: getEnvOrDefault("CURRENT_USER", "fizikurniawan"),
		CurrentTime: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}

	// Database connection (dinamis antara PostgreSQL atau SQLite)
	dbURL := getEnvOrDefault("DATABASE_URL", "sqlite://url_shortener.db")
	var db *sql.DB

	if strings.HasPrefix(dbURL, "postgres://") {
		db, err = sql.Open("postgres", dbURL)
		log.Println("Connected to PostgreSQL")
	} else if strings.HasPrefix(dbURL, "sqlite://") {
		dbPath := strings.TrimPrefix(dbURL, "sqlite://")
		db, err = sql.Open("sqlite", dbPath)
		log.Println("Connected to SQLite")
	} else {
		log.Fatal("Unsupported database URL format")
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Setup session store
	sessionKey := []byte(getEnvOrDefault("SESSION_KEY", "your-secret-key-here-make-it-strong-32"))
	store := sessions.NewCookieStore(sessionKey)

	// Configure session store
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	return &Config{
		DB:    db,
		Store: store,
		Env:   env,
	}
}

func (c *Config) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
