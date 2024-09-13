package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserServiceServer) SwitchBranch(ctx context.Context, req *userProto.SwitchBranchRequest) (*userProto.SwitchBranchResponse, error) {
	res, err := s.userService.SwitchBranch(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
