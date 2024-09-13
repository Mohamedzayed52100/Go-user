package repository

import (
	"context"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/database"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	SharedDbConnection  *gorm.DB
	TenantDBConnections map[string]*gorm.DB
}

func NewDepartmentRepository() *DepartmentRepository {
	postgresService := database.SharedPostgresService

	return &DepartmentRepository{
		SharedDbConnection:  postgresService.Db,
		TenantDBConnections: postgresService.TenantDbConnections,
	}
}

func (r *DepartmentRepository) GetSharedDB() *gorm.DB {
	return r.SharedDbConnection
}

func (r *DepartmentRepository) GetTenantDBConnection(ctx context.Context) *gorm.DB {
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
