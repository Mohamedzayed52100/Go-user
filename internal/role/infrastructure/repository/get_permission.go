package repository

import (
	"context"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"github.com/goplaceapp/goplace-user/utils"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *RoleRepository) GetAllStringPermissions(ctx context.Context) []string {
	var permissions = []string{}

	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil
	}

	if currentUser == nil {
		r.GetTenantDBConnection(ctx).Find(&permissions)
	} else {
		r.GetTenantDBConnection(ctx).
			Model(domain.Permission{}).
			Joins("JOIN role_permission_assignments ON role_permission_assignments.permission_id = permissions.id").
			Where("role_permission_assignments.role_id = ?", currentUser.RoleID).
			Select("permissions.name").
			Scan(&permissions)
	}

	return permissions
}

func (r *RoleRepository) GetAllPermissions(ctx context.Context, req *userProto.GetAllPermissionsRequest) (*userProto.GetAllPermissionsResponse, error) {
	var permissions []*domain.Permission

	r.PermissionsQueryBuilder(ctx, req).
		Group("permissions.id").
		Order("id DESC").
		Find(&permissions)

	return &userProto.GetAllPermissionsResponse{
		Result: utils.CategorizeAndArrangePermissions(permissions),
	}, nil
}

func (r *RoleRepository) PermissionsQueryBuilder(ctx context.Context, req *userProto.GetAllPermissionsRequest) *gorm.DB {
	query := r.GetTenantDBConnection(ctx).
		Model(&domain.Permission{}).
		Joins("LEFT JOIN role_permission_assignments ON role_permission_assignments.permission_id = permissions.id").
		Joins("LEFT JOIN roles ON roles.id = role_permission_assignments.role_id")

	if len(req.GetDepartment()) > 0 {
		departments := []int32{}

		r.GetTenantDBConnection(ctx).
			Model(&departmentDomain.UserDepartment{}).
			Where("id IN ?", req.GetDepartment()).
			Distinct().Pluck("id", &departments)

		query = query.Where("roles.department_id IN (?)", departments)
	}

	if req.GetQuery() != "" {
		searchQuery := "%" + req.GetQuery() + "%"
		query = query.Where("(UPPER(roles.name) LIKE UPPER(?) OR UPPER(roles.display_name) LIKE UPPER(?) OR UPPER(permissions.name) LIKE UPPER(?) OR UPPER(permissions.display_name) LIKE UPPER(?) OR UPPER(permissions.category) LIKE UPPER(?))", searchQuery, searchQuery, searchQuery, searchQuery, searchQuery)
	}
	return query
}

func (r *RoleRepository) GetRolePermissions(ctx context.Context, req *userProto.GetRolePermissionsRequest) (*userProto.GetRolePermissionsResponse, error) {
	permissions, err := r.GetRolePermissionsByRoleID(ctx, req.GetRole())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.GetRolePermissionsResponse{
		Result: permissions,
	}, nil
}

func (r *RoleRepository) GetRolePermissionsByRoleID(ctx context.Context, roleID int32) ([]string, error) {
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
