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
//DTOs

type UserCreateDTO struct{
	Name 	   string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
}


type UserLoginDTO struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateDTO struct{
	Name 	   string `json:"name" binding:"required"`
	Email	   string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type UserResponseDTO struct{
	ID 		  uuid.UUID `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}