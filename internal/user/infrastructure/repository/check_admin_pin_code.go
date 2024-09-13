package repository

import (
	"context"

	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/rbac"
	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
)

func (r *UserRepository) CheckAdminPinCode(ctx context.Context, pinCode string) bool {
	var (
		overbookingPermissionId int32
		users                   = []*externalUserDomain.User{}
	)

	r.SharedDbConnection.
		Model(externalUserDomain.User{}).
		Where("pin_code = ?", pinCode).
		Find(&users)

	r.GetTenantDBConnection(ctx).
		Model(roleDomain.Permission{}).
		Where("name = ?", rbac.OverbookReservation.Name).
		Select("id").
		Scan(&overbookingPermissionId)

	for _, u := range users {
		if err := r.GetTenantDBConnection(ctx).
			First(&roleDomain.RolePermissionAssignment{}, "role_id = ? AND permission_id = ?", u.RoleID, overbookingPermissionId).Error; err != nil {
			return false
		}

	}

	return true
}
