// Package repository handles all direct database operations.
package repository

import (
	"errors"
	"fmt"

	"urlshortener/internal/models"

	"gorm.io/gorm"
)

// URLRepository defines the database operations for URLs.
// Using an interface allows the service layer to swap in a mock during tests.
type URLRepository interface {
	Create(url *models.URL) error
	GetByCode(code string) (*models.URL, error)
	IncrementClicks(id uint) error
	Delete(id uint) error
	List() ([]models.URL, error)
}

type urlRepository struct {
	db *gorm.DB
}

// NewURLRepository returns a PostgreSQL-backed URLRepository.
func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *models.URL) error {
	if err := r.db.Create(url).Error; err != nil {
		return fmt.Errorf("repository: create: %w", err)
	}
	return nil
}

func (r *urlRepository) GetByCode(code string) (*models.URL, error) {
	var record models.URL
	if err := r.db.Where("short_code = ?", code).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repository: get %q: %w", code, ErrNotFound)
		}
		return nil, fmt.Errorf("repository: get: %w", err)
	}
	return &record, nil
}

func (r *urlRepository) IncrementClicks(id uint) error {
	err := r.db.Model(&models.URL{}).
		Where("id = ?", id).
		UpdateColumn("clicks", gorm.Expr("clicks + 1")).Error

	if err != nil {
		return fmt.Errorf("repository: increment clicks: %w", err)
	}
	return nil
}

// Delete soft-deletes the record (sets deleted_at). The row is kept for auditing.
func (r *urlRepository) Delete(id uint) error {
	result := r.db.Delete(&models.URL{}, id)
	if result.Error != nil {
		return fmt.Errorf("repository: delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("repository: delete %d: %w", id, ErrNotFound)
	}
	return nil
}

func (r *urlRepository) List() ([]models.URL, error) {
	var records []models.URL
	if err := r.db.Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("repository: list: %w", err)
	}
	return records, nil
}

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("url not found")
