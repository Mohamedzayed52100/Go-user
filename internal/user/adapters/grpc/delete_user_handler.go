package grpc

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userProto.DeleteUserRequest) (*userProto.DeleteUserResponse, error) {
	res, err := s.userService.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
