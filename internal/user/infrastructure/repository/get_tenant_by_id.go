package repository

import 	"github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"

func (r *UserRepository) GetTenantByID(tenantID string) (*domain.Tenant, error) {
	var tenant *domain.Tenant

	err := r.SharedDbConnection.
		Table("tenants").
		Where("id = ?", tenantID).
		First(&tenant).
		Error
	if err != nil {
		return nil, err
	}

	return tenant, nil
}