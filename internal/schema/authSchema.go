package schema

import "github.com/google/uuid"

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpireIn     string           `json:"expire_in"`
	User         LoginMetaDataDTO `json:"user"`
}

type LoginMetaDataDTO struct {
	ID    uuid.UUID `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}