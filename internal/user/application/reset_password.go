package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) ResetPassword(ctx context.Context, req *userProto.ResetPasswordRequest) (*userProto.ResetPasswordResponse, error) {
	res, err := s.Repository.ResetPassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
