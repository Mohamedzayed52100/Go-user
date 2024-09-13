package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleService) UpdateRole(ctx context.Context, req *userProto.UpdateRoleRequest) (*userProto.UpdateRoleResponse, error) {
	res, err := s.Repository.UpdateRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}