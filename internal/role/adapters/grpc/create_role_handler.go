package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleServiceServer) CreateRole(ctx context.Context, req *userProto.CreateRoleRequest) (*userProto.CreateRoleResponse, error) {
	res, err := s.roleService.CreateRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
