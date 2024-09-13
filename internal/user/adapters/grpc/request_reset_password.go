package grpc

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) RequestResetPassword(ctx context.Context, req *userProto.RequestResetPasswordRequest) (*userProto.RequestResetPasswordResponse, error) {
	res, err := s.userService.RequestResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
