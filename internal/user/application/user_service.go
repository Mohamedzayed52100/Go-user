package application

import "github.com/goplaceapp/goplace-user/internal/user/infrastructure/repository"

type UserService struct {
	Repository *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		Repository: repository.NewUserRepository(),
	}
}
