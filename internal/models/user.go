package models

import (
	"time"

	"github.com/google/uuid"
)



type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name       string
	Email      string
	Password   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string{
	return `"user"`
}
