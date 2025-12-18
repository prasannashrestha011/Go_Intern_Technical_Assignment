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

func ToOrderResponseDTO(order *models.Order)(*schema.OrderResponse){
	return &schema.OrderResponse{
		ID: order.ID,
		UserID: order.UserID,
		Amount: order.Amount,
		CreatedAT: order.CreatedAt,
		Status: order.Status,
	}
}

func ToUserOrderResponseDTO(orders []*models.Order)([]*schema.UserOrderResponse){
	responses:=make([]*schema.UserOrderResponse,0,len(orders))
	for _,order:=range orders{
		if order ==nil{
			continue
		}
	responses = append(responses, &schema.UserOrderResponse{
		ID: order.ID,
		Amount: order.Amount,
		CreatedAt: order.CreatedAt,
		Status: order.Status,
	})
	}
	return responses
}