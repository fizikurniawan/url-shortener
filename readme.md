# URL Shortener Service

A simple URL shortener service built with Go, featuring user authentication, rate limiting, and URL management.

## Features

- User Authentication (Register/Login/Logout)
- URL Shortening
- URL Management Dashboard
- Rate Limiting (5 requests/second with burst of 10)
- Session Management
- Hot Reload for Development

## Tech Stack

- Go 1.21+
- PostgreSQL or Sqlite3
- Tailwind CSS
- gorilla/sessions

## Prerequisites

- Go 1.21 or higher
- PostgreSQL or Sqlite3
- Air (for hot reload)

## Project Structure

```
url-shortener/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── url_handler.go
│   ├── middleware/
│   │   └── auth.go
│   ├── model/
│   │   └── url.go
│   ├── repository/
│   │   └── url_repository.go
│   └── service/
│       └── url_service.go
├── static/
│   └── css/
│       └── style.css
└── templates/
    ├── layout.html
    ├── index.html
    ├── login.html
    └── dashboard.html```
