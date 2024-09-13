package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleServiceServer) DeleteRole(ctx context.Context, req *userProto.DeleteRoleRequest) (*userProto.DeleteRoleResponse, error) {
	res, err := s.roleService.DeleteRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
