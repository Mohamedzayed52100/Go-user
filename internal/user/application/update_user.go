package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) UpdateUser(ctx context.Context, req *userProto.UpdateUserRequest) (*userProto.UpdateUserResponse, error) {
	res, err := s.Repository.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}