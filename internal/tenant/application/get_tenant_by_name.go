package application

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *TenantService) GetTenantByName(ctx context.Context, req *userProto.GetTenantByNameRequest) (*userProto.GetTenantByNameResponse, error) {
	res, err := s.Repository.GetTenantByName(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
