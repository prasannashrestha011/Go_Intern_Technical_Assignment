package repository

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User)(*models.User,error) 
	Delete(id int64) error
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

func (r *userRepo) Delete(id int64) error {
	return r.db.Delete(&models.User{},id).Error
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err:=r.db.Where("email= ?",email).First(&user).Error;err!=nil{
		return nil,err
	}
	return &user,nil
}

func (r *userRepo) GetByID(id int64) (*models.User, error) {
	var user models.User 
	if err:=r.db.Where("id= ?",id).First(&user).Error;err!=nil{
		return nil,err
	}
	return &user,nil
}

func (r *userRepo) Update(user *models.User) (*models.User,error) {
	if err:=r.db.Save(&user).Error;err!=nil{
		return nil,err
	}
	return user,nil
}
