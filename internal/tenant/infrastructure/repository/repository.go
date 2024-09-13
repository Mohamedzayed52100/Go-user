package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/database"
	"gorm.io/gorm"
)

type TenantRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
}

func NewTenantRepository() *TenantRepository {
	postgresService := database.SharedPostgresService

	return &TenantRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
	}
}

func (r *TenantRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *TenantRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
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
