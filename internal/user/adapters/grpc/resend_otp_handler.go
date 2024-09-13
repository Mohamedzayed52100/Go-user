package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) ResendOtp(ctx context.Context, req *userProto.ResendOtpRequest) (*userProto.ResendOtpResponse, error) {
	res, err := s.userService.ResendOtp(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
