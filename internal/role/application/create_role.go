package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleService) CreateRole(ctx context.Context, req *userProto.CreateRoleRequest) (*userProto.CreateRoleResponse, error) {
	res, err := s.Repository.CreateRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
