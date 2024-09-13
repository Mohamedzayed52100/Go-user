package application

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) RequestResetPassword(ctx context.Context, req *userProto.RequestResetPasswordRequest) (*userProto.RequestResetPasswordResponse, error) {
	res, err := s.Repository.RequestResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
