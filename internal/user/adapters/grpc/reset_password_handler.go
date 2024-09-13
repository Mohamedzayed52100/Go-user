package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) ResetPassword(ctx context.Context, req *userProto.ResetPasswordRequest) (*userProto.ResetPasswordResponse, error) {
	res, err := s.userService.ResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
