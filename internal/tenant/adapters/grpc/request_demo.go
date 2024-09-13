package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *TenantServiceServer) RequestDemo(ctx context.Context, req *userProto.DemoRequest) (*userProto.DemoResponse, error) {
	res, err := s.tenantService.RequestDemo(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
