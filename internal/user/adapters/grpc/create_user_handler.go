package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.CreateUserResponse, error) {
	res, err := s.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
