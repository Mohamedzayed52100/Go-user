package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *RoleService) GetAllRoles(ctx context.Context, req *userProto.GetAllRolesRequest) (*userProto.GetAllRolesResponse, error) {
	res, err := s.Repository.GetAllRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, req *userProto.GetRoleByIDRequest) (*userProto.GetRoleByIDResponse, error) {
	res, err := s.Repository.GetRoleByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
