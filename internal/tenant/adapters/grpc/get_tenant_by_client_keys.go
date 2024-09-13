package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *TenantServiceServer) GetTenantByClientKeys(ctx context.Context, req *userProto.GetTenantByClientKeysRequest) (*userProto.GetTenantByClientKeysResponse, error) {
	res, err := s.tenantService.GetTenantByClientKeys(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
