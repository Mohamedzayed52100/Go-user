package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleServiceServer) GetAllRoles(ctx context.Context, req *userProto.GetAllRolesRequest) (*userProto.GetAllRolesResponse, error) {
	res, err := s.roleService.GetAllRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *RoleServiceServer) GetRoleByID(ctx context.Context, req *userProto.GetRoleByIDRequest) (*userProto.GetRoleByIDResponse, error) {
	res, err := s.roleService.GetRoleByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
