package app_grpc

import (
	"context"
	"main/internal/config/protoc"
	"main/internal/services"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserGrpcService struct {
	protoc.UnimplementedUserMicroServiceServer
	UserService services.UserService
}

func (u *UserGrpcService) GetUser(ctx context.Context, req *protoc.UserMicroRequest) (*protoc.UserMicroResponse, error) {

	userDetails, err := u.UserService.GetUserByID(ctx, uuid.MustParse(req.UserId))
	if err != nil {
		return nil, err
	}

	return &protoc.UserMicroResponse{
		Name:       userDetails.Name,
		Email:      userDetails.Email,
		IsVerified: true,
		CreatedAt:  timestamppb.New(userDetails.CreatedAt),
	}, nil
}
