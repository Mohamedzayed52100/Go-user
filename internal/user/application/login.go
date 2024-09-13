package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) Login(ctx context.Context, req *userProto.LoginRequest) (*userProto.LoginResponse, error) {
	///secuiry layer
	res, err := s.Repository.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
