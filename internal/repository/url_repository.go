package repository

import (
	"database/sql"
	"errors"
	"time"
	"url-shortener/internal/model"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(url *model.URL) error {
	query := `
		INSERT INTO urls (long_url, short_code, user_id, visits, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.db.QueryRow(
		query,
		url.LongURL,
		url.ShortCode,
		url.UserID,
		0,
		time.Now(),
	).Scan(&url.ID)
}

func (r *URLRepository) FindByShortCode(code string) (*model.URL, error) {
	url := &model.URL{}
	query := `
		SELECT id, long_url, short_code, user_id, visits, created_at
		FROM urls
		WHERE short_code = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, code).Scan(
		&url.ID,
		&url.LongURL,
		&url.ShortCode,
		&url.UserID,
		&url.Visits,
		&url.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("url not found")
	}

	return url, err
}

func (r *URLRepository) GetUserURLs(userID int64) ([]model.URL, error) {
	var urls []model.URL
	query := `
		SELECT id, long_url, short_code, user_id, visits, created_at, deleted_at
		FROM urls
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url model.URL
		err := rows.Scan(
			&url.ID,
			&url.LongURL,
			&url.ShortCode,
			&url.UserID,
			&url.Visits,
			&url.CreatedAt,
			&url.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (r *URLRepository) IncrementVisits(id int64) error {
	query := `
		UPDATE urls
		SET visits = visits + 1
		WHERE id = $1`

	_, err := r.db.Exec(query, id)
	return err
}

func (r *URLRepository) IsShortCodeAvailable(code string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	r.db.QueryRow(query, code).Scan(&exists)
	return !exists
}
