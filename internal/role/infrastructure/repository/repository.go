package repository

import (
	"context"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/database"
	"gorm.io/gorm"
)

type RoleRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
}

func NewRoleRepository() *RoleRepository {
	postgresService := database.SharedPostgresService

	return &RoleRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
	}
}

func (r *RoleRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *RoleRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
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
