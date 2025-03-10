package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength = 4 // Default length untuk random code
)

func (s *URLService) CreateShortURL(longURL string, userID *int64, customCode string) (*model.URL, error) {
	// Validasi URL
	if strings.TrimSpace(longURL) == "" {
		return nil, errors.New("URL cannot be empty")
	}

	url := &model.URL{
		LongURL: longURL,
		UserID:  userID,
	}

	// Handle custom code untuk user yang login
	if customCode != "" {
		if userID == nil {
			return nil, errors.New("custom codes are only available for logged in users")
		}

		// Validasi custom code
		if len(customCode) < 3 {
			return nil, errors.New("custom code must be at least 3 characters long")
		}

		if !isValidCustomCode(customCode) {
			return nil, errors.New("custom code can only contain letters, numbers, and hyphens")
		}

		// Cek availability
		if !s.repo.IsShortCodeAvailable(customCode) {
			return nil, errors.New("this custom code is already taken")
		}

		url.ShortCode = customCode
	} else {
		// Generate random code untuk anonymous user
		code, err := s.generateUniqueCode()
		if err != nil {
			return nil, err
		}
		url.ShortCode = code
	}

	err := s.repo.Create(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) generateUniqueCode() (string, error) {
	length := codeLength
	maxAttempts := 3 // Jumlah percobaan sebelum menambah panjang

	for attempt := 0; attempt < maxAttempts; attempt++ {
		code, err := s.generateRandomString(length)
		if err != nil {
			return "", err
		}

		if s.repo.IsShortCodeAvailable(code) {
			return code, nil
		}
	}

	// Jika gagal dengan panjang default, coba dengan panjang yang lebih besar
	length++
	code, err := s.generateRandomString(length)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *URLService) generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func isValidCustomCode(code string) bool {
	for _, char := range code {
		if !strings.ContainsRune(charset+"-", char) {
			return false
		}
	}
	return true
}

func (s *URLService) GetURL(shortCode string) (*model.URL, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	err = s.repo.IncrementVisits(url.ID)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetUserURLs(userID int64) ([]model.URL, error) {
	return s.repo.GetUserURLs(userID)
}
