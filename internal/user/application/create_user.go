package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.CreateUserResponse, error) {
	res, err := s.Repository.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}