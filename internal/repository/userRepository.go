package repository

import (
	"main/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uuid.UUID) error
}

type userRepo struct {
	db *gorm.DB
}


func NewRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

// pagination will be implemented in the future
func (r *userRepo) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err:=r.db.Find(users).Error;err!=nil{
		return nil,err
	}
	return users,nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email= ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id= ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
