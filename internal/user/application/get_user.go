package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) GetAllUsers(ctx context.Context, req *userProto.GetAllUsersRequest) (*userProto.GetAllUsersResponse, error) {
	res, err := s.Repository.GetAllUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req *userProto.GetUserByIDRequest) (*userProto.GetUserByIDResponse, error) {
	res, err := s.Repository.GetUserByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserService) GetAuthenticatedUser(ctx context.Context, req *emptypb.Empty) (*userProto.GetAuthenticatedUserResponse, error) {
	res, err := s.Repository.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
