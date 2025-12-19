package repository

import (
	"context"
	"errors"
	"main/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context,user *models.User) error
	GetByID(ctx context.Context,id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context,email string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context,user *models.User) (*models.User, error)
	Delete(ctx context.Context,id uuid.UUID) error
}

type userRepo struct {
	db *gorm.DB
}


func NewRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context,user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) Delete(ctx context.Context,id uuid.UUID) error {
	result:= r.db.WithContext(ctx).Delete(&models.User{},"id= ?", id)
	if result.Error !=nil{
		return result.Error
	}
	if result.RowsAffected==0{
		return errors.New("User details not found") 
	}
	return nil

}

// pagination will be implemented in the future
func (r *userRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err:=r.db.WithContext(ctx).Find(&users).Error;err!=nil{
		return nil,err
	}
	return users,nil
}

func (r *userRepo) GetByEmail(ctx context.Context,email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email= ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(ctx context.Context,id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id= ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context,user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
