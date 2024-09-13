package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleServiceServer) UpdateRole(ctx context.Context, req *userProto.UpdateRoleRequest) (*userProto.UpdateRoleResponse, error) {
	res, err := s.roleService.UpdateRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
