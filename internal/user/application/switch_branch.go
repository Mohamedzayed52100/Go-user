package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *UserService) SwitchBranch(ctx context.Context, req *userProto.SwitchBranchRequest) (*userProto.SwitchBranchResponse, error) {
	res, err := s.Repository.SwitchBranch(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
