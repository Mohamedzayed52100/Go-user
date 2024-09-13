package grpc

import "github.com/goplaceapp/goplace-user/internal/user/application"

type UserServiceServer struct {
	userService *application.UserService
}

func NewUserService() *UserServiceServer {
	return &UserServiceServer{
		userService: application.NewUserService(),
	}
}
