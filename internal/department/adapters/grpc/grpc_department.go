package grpc

import "github.com/goplaceapp/goplace-user/internal/department/application"

type DepartmentServiceServer struct {
	departmentService *application.DepartmentService
}

func NewDepartmentService() *DepartmentServiceServer {
	return &DepartmentServiceServer{
		departmentService: application.NewDepartmentService(),
	}
}
