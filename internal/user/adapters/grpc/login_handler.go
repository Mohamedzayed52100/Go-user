package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) Login(ctx context.Context, req *userProto.LoginRequest) (*userProto.LoginResponse, error) {
	res, err := s.userService.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
