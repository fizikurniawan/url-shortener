package model

import "time"

type URL struct {
	ID        int64      `json:"id"`
	LongURL   string     `json:"long_url"`
	ShortCode string     `json:"short_code"`
	UserID    *int64     `json:"user_id"`
	Visits    int64      `json:"visits"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsActive bool   `json:"is_active"`
}
