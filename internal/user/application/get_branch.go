package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) GetAllBranches(ctx context.Context, req *emptypb.Empty) (*userProto.GetAllBranchesResponse, error) {
	res, err := s.Repository.GetAllBranches(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
