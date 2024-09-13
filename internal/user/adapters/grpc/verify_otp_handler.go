package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) VerifyOtp(ctx context.Context, req *userProto.VerifyOtpRequest) (*userProto.VerifyOtpResponse, error) {
	res, err := s.userService.VerifyOtp(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
