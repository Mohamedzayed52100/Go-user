package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) ChangePassword(ctx context.Context, req *userProto.ChangePasswordRequest) (*userProto.ChangePasswordResponse, error) {
	res, err := s.userService.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
