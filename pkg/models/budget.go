package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	gorm.Model
	Name   string `gorm:"size:255;not null" json:"name"`
	Amount int64  `gorm:"not null" json:"amount"`
	UserID uuid.UUID
}
