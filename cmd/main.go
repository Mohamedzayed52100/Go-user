package main

import (
	"github.com/goplaceapp/goplace-common/pkg/grpchelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/server"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"os"
)

func main() {
	godotenv.Load()
	meta.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	service := server.New()
	grpchelper.Start(func(server *grpc.Server) {
		userProto.RegisterUserServer(server, service.UserServiceServer)
		userProto.RegisterRoleServer(server, service.RoleServiceServer)
		userProto.RegisterDepartmentServer(server, service.DepartmentServiceServer)
		userProto.RegisterTenantServer(server, service.TenantServiceServer)
	}, service)
}
