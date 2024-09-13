package repository

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/tenant/adapters/convertors"
	"github.com/goplaceapp/goplace-user/internal/tenant/domain"
	tenantDomain "github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"

)

func (r *TenantRepository) GetTenantByClientKeys(ctx context.Context, req *userProto.GetTenantByClientKeysRequest) (*userProto.GetTenantByClientKeysResponse, error) {
	var credentials domain.TenantCredential

	err := r.SharedDbConnection.
		Model(&domain.TenantCredential{}).
		Where("client_id = ? AND client_secret = ?", req.GetClientId(), req.GetClientSecret()).
		First(&credentials).Error
	if err != nil {
		return nil, err
	}

	var tenant tenantDomain.Tenant
	err = r.SharedDbConnection.
		Model(&tenantDomain.Tenant{}).
		Where("id = ?", credentials.TenantID).
		First(&tenant).Error
	if err != nil {
		return nil, err
	}

	return &userProto.GetTenantByClientKeysResponse{
		Result: convertors.BuildTenantResponse(&tenant),
	}, nil
}
