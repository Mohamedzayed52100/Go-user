package repository

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/tenant/adapters/convertors"
	"github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"
)

func (r *TenantRepository) GetTenantByName(ctx context.Context, req *userProto.GetTenantByNameRequest) (*userProto.GetTenantByNameResponse, error) {
	var tenant domain.Tenant

	err := r.SharedDbConnection.
		Model(&domain.Tenant{}).
		Joins("LEFT JOIN tenant_profiles ON tenant_profiles.tenant_id = tenants.id").
		Where("tenant_profiles.name = ?", req.GetName()).
		First(&tenant).Error
	if err != nil {
		return nil, err
	}

	return &userProto.GetTenantByNameResponse{
		Result: convertors.BuildTenantResponse(&tenant),
	}, nil
}
