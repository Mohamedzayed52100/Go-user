package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) ChangePassword(ctx context.Context, req *userProto.ChangePasswordRequest) (*userProto.ChangePasswordResponse, error) {
	res, err := s.Repository.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
