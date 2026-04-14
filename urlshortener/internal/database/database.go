package database

import (
	"fmt"
	"log"
	"time"

	"urlshortener/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a PostgreSQL connection and runs migrations.
func Connect(dsn string, isDev bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if isDev {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: true, // cache compiled statements for faster repeated queries
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("database connected")

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// runMigrations creates or updates tables and handles column consolidations.
func runMigrations(db *gorm.DB) error {
	// If the old 'original' column exists, try to migrate it to 'original_url'
	_ = db.Exec(`
		DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='urls' AND column_name='original') THEN
				IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='urls' AND column_name='original_url') THEN
					-- Both exist: migrate data
					UPDATE urls SET original_url = original WHERE original_url IS NULL OR original_url = '';
					-- Make original nullable so it doesn't block inserts
					ALTER TABLE urls ALTER COLUMN original DROP NOT NULL;
				ELSE
					-- Only original exists: rename it
					ALTER TABLE urls RENAME COLUMN original TO original_url;
				END IF;
			END IF;
		END $$;
	`)

	if err := db.AutoMigrate(&models.URL{}); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("database: migrations applied")
	return nil
}
