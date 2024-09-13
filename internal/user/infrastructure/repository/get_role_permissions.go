package repository

import (
	"context"

	"github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
)

func (r *UserRepository) GetRolePermissionsByRoleID(ctx context.Context, roleID int32) ([]string, error) {
	permissions := []*domain.Permission{}

	if err := r.GetTenantDBConnection(ctx).Model(&permissions).
		Joins("JOIN role_permission_assignments ON role_permission_assignments.permission_id = permissions.id").
		Where("role_permission_assignments.role_id = ?", roleID).
		Select("permissions.*").
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	stringPermissions := []string{}
	for _, p := range permissions {
		stringPermissions = append(stringPermissions, p.Name)
	}

	return stringPermissions, nil
}
