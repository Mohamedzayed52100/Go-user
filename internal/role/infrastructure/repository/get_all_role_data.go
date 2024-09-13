package repository

import (
	"context"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	"net/http"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"google.golang.org/grpc/status"
)

func (r *RoleRepository) GetAllRoleData(ctx context.Context, id int32) (*externalRoleDomain.Role, error) {
	var role *externalRoleDomain.Role
	if err := r.GetTenantDBConnection(ctx).
		Where("id =?", id).
		First(&role).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var department *departmentDomain.UserDepartment
	if err := r.GetTenantDBConnection(ctx).
		Model(&departmentDomain.UserDepartment{}).
		Where("id = ?", role.DepartmentID).
		First(&department).
		Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	role.Department = department

	permissions := []*roleDomain.Permission{}
	r.GetTenantDBConnection(ctx).
		Model(&roleDomain.Permission{}).
		Joins("JOIN role_permission_assignments ON role_permission_assignments.permission_id = permissions.id").
		Where("role_id =?", role.ID).
		Find(&permissions)
	role.Permissions = permissions

	// Count all users assigned to this role
	var count int64
	r.GetSharedDB().
		Model(&externalUserDomain.User{}).
		Where("role_id = ?", role.ID).
		Count(&count)
	role.UsersCount = int32(count)

	return role, nil
}
