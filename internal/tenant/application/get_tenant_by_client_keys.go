package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *TenantService) GetTenantByClientKeys(ctx context.Context, req *userProto.GetTenantByClientKeysRequest) (*userProto.GetTenantByClientKeysResponse, error) {
	res, err := s.Repository.GetTenantByClientKeys(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
