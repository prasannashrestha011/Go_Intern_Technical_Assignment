package services

import (
	"context"
	"main/internal/models"
	"main/internal/repository"
	"main/internal/schema"
	"main/internal/utils"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, newOrder *schema.CreateOrder) error
	GetOrders(ctx context.Context) ([]*models.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*schema.OrderResponse, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*schema.UserOrderResponse, error)
}

type orderService struct {
	repo repository.OrderRepository
}


func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, newOrder *schema.CreateOrder) error {
	newOrderModel:=&models.Order{
		OrderName: newOrder.OrderName,
		UserID: newOrder.UserID,
		Amount: newOrder.Amount,
		Status: "pending",
	}
	return o.repo.Create(ctx,newOrderModel)
	
}

func (o *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*schema.OrderResponse, error) {
	orderModel,err:=o.repo.Get(ctx,id)
	if err!=nil{
		return nil,err
	}
	orderDTO:=utils.ToOrderResponseDTO(orderModel)
	return  orderDTO,nil
}

func (o *orderService) GetOrders(ctx context.Context) ([]*models.Order, error) {
	orders,err:=o.repo.GetAll(ctx)
	if err!=nil{
		return nil,err
	}
	return orders,nil
}

func (o *orderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*schema.UserOrderResponse, error) {

	userOrders,err:=o.repo.GetUserOrders(ctx,userID)
	if err!=nil{
		return nil,err
	}
	userOrdersDTOs:=utils.ToUserOrderResponseDTO(userOrders)

	return userOrdersDTOs,nil 
}
