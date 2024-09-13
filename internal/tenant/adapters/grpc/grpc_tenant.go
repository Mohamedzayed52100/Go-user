package grpc

import "github.com/goplaceapp/goplace-user/internal/tenant/application"

type TenantServiceServer struct {
	tenantService *application.TenantService
}

func NewTenantService() *TenantServiceServer {
	return &TenantServiceServer{
		tenantService: application.NewTenantService(),
	}
}
