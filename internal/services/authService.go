package services

import (
	"context"
	"main/internal/logger"
	"main/internal/repository"
	"main/internal/schema"
	"main/internal/utils"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, creds *schema.UserLoginDTO) (*schema.LoginMetaDataDTO, error)
}

type authService struct {
	repo repository.UserRepository
}


func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}


func (a *authService) Login(ctx context.Context, creds *schema.UserLoginDTO) (*schema.LoginMetaDataDTO, error) {
	email:=creds.Email

	userDetails,err:=a.repo.GetByEmail(ctx,email)

	if err!=nil{
		logger.Log.Error("Failed to fetch the user details",zap.Error(err))
		return nil,err
	}
	isMatches:=utils.ComparePassword(userDetails.Password,creds.Password)

	if !isMatches{
		logger.Log.Error("Authentication Error: ",zap.Error(err))
		return nil,err
	}

	userMetaData:=&schema.LoginMetaDataDTO{
		ID: userDetails.ID,
		Name: userDetails.Name,
		Email: userDetails.Email,
	}

	return userMetaData,nil
}