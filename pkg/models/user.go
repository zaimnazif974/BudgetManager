package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// uuid.Nil adalah UUID "zero"
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return
}

type User struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email      string    `gorm:"uniqueIndex"`
	Provider   string
	Password   string
	FirstName  string
	LastName   string
	AcessToken string
}
