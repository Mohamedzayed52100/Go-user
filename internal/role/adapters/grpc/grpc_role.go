package grpc

import "github.com/goplaceapp/goplace-user/internal/role/application"

type RoleServiceServer struct {
	roleService *application.RoleService
}

func NewRoleService() *RoleServiceServer {
	return &RoleServiceServer{
		roleService: application.NewRoleService(),
	}
}
