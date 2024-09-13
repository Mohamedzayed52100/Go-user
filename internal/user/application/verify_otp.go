package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) VerifyOtp(ctx context.Context, req *userProto.VerifyOtpRequest) (*userProto.VerifyOtpResponse, error) {
	res, err := s.Repository.VerifyOtp(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
