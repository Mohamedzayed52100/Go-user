package grpc

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/goplaceapp/goplace-user/config"
	"github.com/goplaceapp/goplace-user/database"
	"github.com/goplaceapp/goplace-user/internal/user/application"
	"gorm.io/gorm"
)

type UserServiceServer struct {
	UserService *application.UserService
}

func NewUserService(db *gorm.DB, tenantDBConnections map[string]*gorm.DB) *UserServiceServer {
	if database.SharedPostgresService == nil {
		cfg := &config.Config{}
		if err := env.Parse(cfg); err != nil {
			panic(fmt.Errorf("failed to parse environment variables, %w", err))
		}

		database.SharedPostgresService = &database.PostgresService{
			Db:                  db,
			TenantDbConnections: tenantDBConnections,
			SvcCfg:              cfg,
		}
	}

	return &UserServiceServer{
		UserService: application.NewUserService(),
	}
}
