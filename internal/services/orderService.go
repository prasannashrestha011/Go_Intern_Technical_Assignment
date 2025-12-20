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
	CreateOrder(ctx context.Context, newOrder *schema.CreateOrder)(*schema.OrderResponse,error) 
	GetOrders(ctx context.Context,page int,pageSize int) ([]*schema.OrderResponse, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*schema.OrderResponse, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID,page int,pageSize int) ([]*schema.UserOrderResponse, error)
	UpdateOrderDetails(ctx context.Context,id uuid.UUID, order *schema.OrderUpdate)(*schema.OrderResponse,error)
}

type orderService struct {
	repo repository.OrderRepository
}


func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, newOrder *schema.CreateOrder) (*schema.OrderResponse,error) {
	newOrderModel := &models.Order{
		OrderName: newOrder.OrderName,
		UserID:    newOrder.UserID,
		Quantity: newOrder.Quantity,
		Price: newOrder.Price,
		Amount:    newOrder.Price*float64(newOrder.Quantity),
		Status:    "pending",
	}

	order,err:= o.repo.Create(ctx, newOrderModel)
	if err!=nil{
		return nil,err
	}
	utils.OrderQueue<-*order
	orderResp:=utils.ToOrderResponseDTO(order)
	return orderResp,nil

}

func (o *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*schema.OrderResponse, error) {
	orderModel, err := o.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	orderDTO := utils.ToOrderResponseDTO(orderModel)
	return orderDTO, nil
}

func (o *orderService) GetOrders(ctx context.Context,page int, pageSize int) ([]*schema.OrderResponse, error) {
	orders, err := o.repo.GetAll(ctx,page,pageSize)

	if err != nil {
		return nil, err
	}
	ordersDTO:=make([]*schema.OrderResponse,len(orders))
	for i,order:=range orders{
		ordersDTO[i] = utils.ToOrderResponseDTO(order)
	}
	return ordersDTO, nil
}

func (o *orderService) GetUserOrders(ctx context.Context, userID uuid.UUID,page int,pageSize int) ([]*schema.UserOrderResponse, error) {


	userOrders, err := o.repo.GetUserOrders(ctx, userID,page,pageSize)
	if err != nil {
		return nil, err
	}
	userOrdersDTOs := utils.ToUserOrderResponseDTO(userOrders)

	return userOrdersDTOs, nil
}

func (o *orderService) UpdateOrderDetails(
	ctx context.Context,
	id uuid.UUID,
	updateDto *schema.OrderUpdate,
) (*schema.OrderResponse, error) {
	existing, err := o.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if updateDto.OrderName != nil {
	existing.OrderName = *updateDto.OrderName
	}
	if updateDto.Price != nil {
		existing.Price = *updateDto.Price
	}
	if updateDto.Quantity !=nil{
		existing.Quantity = *updateDto.Quantity
	}

	existing.Amount=existing.Price*float64(existing.Quantity);
	if updateDto.Status != nil {
		existing.Status = *updateDto.Status
	}

	
	if _,err:=o.repo.Update(ctx,existing);err!=nil{
		return nil,err
	}

	return utils.ToOrderResponseDTO(existing), nil
}
