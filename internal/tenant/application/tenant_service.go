package application

import "github.com/goplaceapp/goplace-user/internal/tenant/infrastructure/repository"

type TenantService struct {
	Repository *repository.TenantRepository
}

func NewTenantService() *TenantService {
	return &TenantService{
		Repository: repository.NewTenantRepository(),
	}
}
