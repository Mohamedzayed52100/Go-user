package application

import "github.com/goplaceapp/goplace-user/internal/department/infrastructure/repository"

type DepartmentService struct {
	Repository *repository.DepartmentRepository
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		Repository: repository.NewDepartmentRepository(),
	}
}
