// Package service contains all business logic for the URL shortener.
// It sits between HTTP handlers and the repository.
// No Fiber or HTTP types appear here — keeping it easy to test.
package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"urlshortener/internal/models"
	"urlshortener/internal/repository"
)

// DTOs
type ShortenRequest struct {
	URL            string `json:"url"`
	CustomCode     string `json:"custom_code"`
	ExpiresInHours int    `json:"expires_in_hours"`
}

type URLResponse struct {
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	Clicks      int64      `json:"clicks"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// URLService defines the business logic.
type URLService interface {
	Shorten(req ShortenRequest) (*URLResponse, error)
	GetOriginal(code string) (string, error)
	GetStats(code string) (*URLResponse, error)
	Delete(code string) error
	List() ([]URLResponse, error)
}

type urlService struct {
	repo    repository.URLRepository
	baseURL string
}

// NewURLService returns a new URLService.
func NewURLService(repo repository.URLRepository, baseURL string) URLService {
	return &urlService{
		repo:    repo,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *urlService) Shorten(req ShortenRequest) (*URLResponse, error) {
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}

	code := req.CustomCode
	if code == "" {
		generated, err := generateCode(6)
		if err != nil {
			return nil, fmt.Errorf("service: generate code: %w", err)
		}
		code = generated
	} else if err := validateCode(code); err != nil {
		return nil, err
	}

	newURL := &models.URL{
		ShortCode:   code,
		OriginalURL: req.URL,
	}

	if req.ExpiresInHours > 0 {
		expiry := time.Now().Add(time.Duration(req.ExpiresInHours) * time.Hour)
		newURL.ExpiresAt = &expiry
	}

	if err := s.repo.Create(newURL); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, ErrCodeTaken
		}
		return nil, err
	}

	return s.mapToResponse(newURL), nil
}

func (s *urlService) GetOriginal(code string) (string, error) {
	record, err := s.repo.GetByCode(code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	if record.IsExpired() {
		return "", ErrExpired
	}

	go func() { _ = s.repo.IncrementClicks(record.ID) }()

	return record.OriginalURL, nil
}

func (s *urlService) GetStats(code string) (*URLResponse, error) {
	record, err := s.repo.GetByCode(code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.mapToResponse(record), nil
}

func (s *urlService) Delete(code string) error {
	record, err := s.repo.GetByCode(code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(record.ID)
}

func (s *urlService) List() ([]URLResponse, error) {
	records, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	res := make([]URLResponse, len(records))
	for i := range records {
		res[i] = *s.mapToResponse(&records[i])
	}
	return res, nil
}

// Helpers
func (s *urlService) mapToResponse(u *models.URL) *URLResponse {
	return &URLResponse{
		ShortCode:   u.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, u.ShortCode),
		OriginalURL: u.OriginalURL,
		Clicks:      u.Clicks,
		ExpiresAt:   u.ExpiresAt,
		CreatedAt:   u.CreatedAt,
	}
}

func generateCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func validateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("invalid URL: must be absolute http/https")
	}
	return nil
}

func validateCode(code string) error {
	if len(code) < 3 || len(code) > 20 {
		return fmt.Errorf("code must be 3-20 characters")
	}
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("code includes invalid characters")
		}
	}
	return nil
}

// Errors
var (
	ErrNotFound  = errors.New("url not found")
	ErrExpired   = errors.New("url expired")
	ErrCodeTaken = errors.New("short code taken")
)
