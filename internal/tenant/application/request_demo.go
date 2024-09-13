package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *TenantService) RequestDemo(ctx context.Context, req *userProto.DemoRequest) (*userProto.DemoResponse, error) {
	res, err := s.Repository.RequestDemo(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
