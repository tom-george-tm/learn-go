package models

import (
	"time"

	"gorm.io/gorm"
)

// URL maps to the "urls" table in PostgreSQL.
type URL struct {
	gorm.Model
	ShortCode   string     `gorm:"uniqueIndex;not null;size:20" json:"short_code"`
	OriginalURL string     `gorm:"not null;type:text"          json:"original_url"`
	Clicks      int64      `gorm:"default:0"                   json:"clicks"`
	ExpiresAt   *time.Time `gorm:"index"                       json:"expires_at,omitempty"`
}

// IsExpired returns true if the URL has passed its expiry time.
func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}
