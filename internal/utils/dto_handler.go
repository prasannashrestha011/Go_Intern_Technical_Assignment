// Methods to convert database model into DTO objects for response schema
package utils

import (
	"main/internal/models"
	"main/internal/schema"
)


func ToUserResponseDTO(user *models.User)(*schema.UserResponseDTO){
	return &schema.UserResponseDTO{
				ID: user.ID,
				Name: user.Name,
				Email: user.Email,
				CreatedAt: user.CreatedAt,
			}
}