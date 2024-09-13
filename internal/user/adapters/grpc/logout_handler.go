package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserServiceServer) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	res, err := s.userService.Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
