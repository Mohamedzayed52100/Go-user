package application

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	res, err := s.Repository.Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
