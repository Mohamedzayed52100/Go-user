package application

import "github.com/goplaceapp/goplace-user/internal/role/infrastructure/repository"

type RoleService struct {
	Repository *repository.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		Repository: repository.NewRoleRepository(),
	}
}
