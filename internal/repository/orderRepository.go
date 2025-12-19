package repository

import (
	"context"
	"fmt"
	"main/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, newOrder *models.Order) error
	Get(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	GetAll(ctx context.Context) ([]*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}


func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}

}
func (o *orderRepository) Create(ctx context.Context, newOrder *models.Order) error {
	return o.db.WithContext(ctx).Create(newOrder).Error
}

func (o *orderRepository) Get(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var fetched_order models.Order
	err := o.db.WithContext(ctx).Where("id= ?", id).First(&fetched_order).Error
	if err != nil {
		return nil, err
	}
	return &fetched_order, nil
}

func (o *orderRepository) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	var user_orders []*models.Order
	err := o.db.WithContext(ctx).Where("user_id= ?", userID).Find(&user_orders).Error
	if err != nil {
		return nil, err
	}
    if len(user_orders) == 0 {
        return nil, fmt.Errorf("no orders found for user %s", userID)
    }
	return user_orders, nil
}

func (o *orderRepository) GetAll(ctx context.Context) ( orders []*models.Order, err error) {
	err=o.db.WithContext(ctx).Find(&orders).Error
	if err!=nil{
		return nil,err
	}
	return orders,nil
}