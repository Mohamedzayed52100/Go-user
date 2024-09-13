package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleService) GetAllPermissions(ctx context.Context, req *userProto.GetAllPermissionsRequest) (*userProto.GetAllPermissionsResponse,error) {
	res, err := s.Repository.GetAllPermissions(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *RoleService) GetRolePermissions(ctx context.Context, req *userProto.GetRolePermissionsRequest) (*userProto.GetRolePermissionsResponse, error) {
	res, err := s.Repository.GetRolePermissions(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
