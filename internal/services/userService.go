package services

import (
	"context"
	"log"
	"main/internal/models"
	"main/internal/repository"
	"main/internal/utils"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context,user *models.UserCreateDTO) (*models.UserResponseDTO, error)
	GetUserByID(ctx context.Context,id uuid.UUID) (*models.UserResponseDTO, error)
	GetUsers(ctx context.Context) ([]*models.UserResponseDTO, error)
	GetUserByEmail(ctx context.Context,email string) (*models.UserResponseDTO, error)
	UpdateUser(ctx context.Context,id uuid.UUID, user *models.UserUpdateDTO) (*models.UserResponseDTO, error)
	DeleteUser(ctx context.Context,id uuid.UUID) error
}

type userService struct {
	repo repository.UserRepository
}

func (u *userService) GetUsers(ctx context.Context) ([]*models.UserResponseDTO, error) {
	users,err:=u.repo.GetAll(ctx)
	if err!=nil{
		return nil,err
	}

	var usersDTO []*models.UserResponseDTO
	for _,user:=range users{
		dto:=&models.UserResponseDTO{
			ID:user.ID ,
			Name: user.Name,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		usersDTO = append(usersDTO,dto )
	}
	return usersDTO,nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) RegisterUser(ctx context.Context,dto *models.UserCreateDTO) (*models.UserResponseDTO, error) {
	
	hashedPwd,err:=utils.HashPassoword(dto.Password)
	if err!=nil{
		return nil,err
	}
	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password:hashedPwd,
	}
	err = u.repo.Create(ctx,user)
	if err != nil {
		return nil, err
	}
	userDTO := utils.ToUserResponseDTO(user)
	return userDTO, nil
}

func (u *userService) DeleteUser(ctx context.Context,id uuid.UUID) error {
	return u.repo.Delete(ctx,id)
}

func (u *userService) UpdateUser(ctx context.Context,id uuid.UUID, dto *models.UserUpdateDTO) (*models.UserResponseDTO, error) {

	user, err := u.repo.GetByID(ctx,id)
	if err != nil {
		return nil, err
	}
	log.Println(dto)
	if dto.Name != "" {
		user.Name = dto.Name
	}

	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Password != "" {
		user.Password = dto.Password
	}

	updated_user, err := u.repo.Update(ctx,user)
	if err != nil {
		return nil, err
	}
	userDTO := utils.ToUserResponseDTO(updated_user)
	return userDTO, nil
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*models.UserResponseDTO, error) {
	user, err := u.repo.GetByEmail(ctx,email)
	if err != nil {
		return nil, err
	}
	userDTO := utils.ToUserResponseDTO(user)
	return userDTO, nil
}

func (u *userService) GetUserByID(ctx context.Context,id uuid.UUID) (*models.UserResponseDTO, error) {
	user, err := u.repo.GetByID(ctx,id)
	if err != nil {
		return nil, err
	}
	userDTO := utils.ToUserResponseDTO(user)
	return userDTO, nil
}
