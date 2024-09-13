package grpc

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *TenantServiceServer) GetTenantByName(ctx context.Context, req *userProto.GetTenantByNameRequest) (*userProto.GetTenantByNameResponse, error) {
	res, err := s.tenantService.GetTenantByName(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
