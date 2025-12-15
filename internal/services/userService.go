package services

import (
	"main/internal/models"
	"main/internal/repository"
)

type UserService interface {
	RegisterUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserRepository(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
// DeleteUser implements [UserService].
func (u *userService) DeleteUser(id int64) error {
	panic("unimplemented")
}

// GetUserByEmail implements [UserService].
func (u *userService) GetUserByEmail(email string) (*models.User, error) {
	panic("unimplemented")
}

// GetUserByID implements [UserService].
func (u *userService) GetUserByID(id int64) (*models.User, error) {
	panic("unimplemented")
}

// RegisterUser implements [UserService].
func (u *userService) RegisterUser(user *models.User) error {
	panic("unimplemented")
}

// UpdateUser implements [UserService].
func (u *userService) UpdateUser(user *models.User) (*models.User, error) {
	panic("unimplemented")
}

