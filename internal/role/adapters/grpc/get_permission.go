package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleServiceServer) GetAllPermissions(ctx context.Context, req *userProto.GetAllPermissionsRequest) (*userProto.GetAllPermissionsResponse, error) {
	res, err := s.roleService.GetAllPermissions(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *RoleServiceServer) GetRolePermissions(ctx context.Context, req *userProto.GetRolePermissionsRequest) (*userProto.GetRolePermissionsResponse, error) {
	res, err := s.roleService.GetRolePermissions(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
