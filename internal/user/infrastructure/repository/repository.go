package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/database"
	"github.com/goplaceapp/goplace-user/internal/role/infrastructure/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
	RoleRepository      *repository.RoleRepository
}

func NewUserRepository() *UserRepository {
	postgresService := database.SharedPostgresService

	return &UserRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
		RoleRepository:      repository.NewRoleRepository(),
	}
}

func (r *UserRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *UserRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
	tenantDBNameValue := ctx.Value(meta.TenantDBNameContextKey.String())

	if tenantDBNameValue == nil {
		return nil
	}
	tenantDBName, ok := tenantDBNameValue.(string)
	if !ok {
		return nil
	}

	if db, ok := r.TenantDBConnections[tenantDBName]; ok {
		return db
	}

	return nil
}
