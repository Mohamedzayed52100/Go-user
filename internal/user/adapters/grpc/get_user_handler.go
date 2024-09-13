package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserServiceServer) GetAllUsers(ctx context.Context, req *userProto.GetAllUsersRequest) (*userProto.GetAllUsersResponse, error) {
	res, err := s.userService.GetAllUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserServiceServer) GetUserByID(ctx context.Context, req *userProto.GetUserByIDRequest) (*userProto.GetUserByIDResponse, error) {
	res, err := s.userService.GetUserByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserServiceServer) GetAuthenticatedUser(ctx context.Context, req *emptypb.Empty) (*userProto.GetAuthenticatedUserResponse, error) {
	res, err := s.userService.GetAuthenticatedUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
