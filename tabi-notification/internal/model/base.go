package model

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common fields for all models
// Do not use gorm.Model because of uint ID
type Base struct {
	// The time that record is created
	CreatedAt time.Time `json:"created_at"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"updated_at"`
	// The time that record is deleted
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
