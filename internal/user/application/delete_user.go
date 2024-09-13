package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) DeleteUser(ctx context.Context, req *userProto.DeleteUserRequest) (*userProto.DeleteUserResponse, error) {
	res, err := s.Repository.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}