package grpc

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/goplaceapp/goplace-user/config"
	"github.com/goplaceapp/goplace-user/database"
	"github.com/goplaceapp/goplace-user/internal/role/application"
	"gorm.io/gorm"
)

type RoleServiceServer struct {
	RoleService *application.RoleService
}

func NewRoleService(db *gorm.DB, tenantDBConnections map[string]*gorm.DB) *RoleServiceServer {
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

	return &RoleServiceServer{
		RoleService: application.NewRoleService(),
	}
}
