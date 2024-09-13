package repository

import (
	"context"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	"net/http"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/role/adapters/convertors"
	"github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"github.com/goplaceapp/goplace-user/utils"
	"google.golang.org/grpc/status"
)

func (r *RoleRepository) UpdateRole(ctx context.Context, req *userProto.UpdateRoleRequest) (*userProto.UpdateRoleResponse, error) {
	if err := r.GetTenantDBConnection(ctx).
		First(
			&externalRoleDomain.Role{}, "id = ?",
			req.GetParams().GetId(),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Role not found")
	}

	if req.GetParams().GetDepartment() != 0 {
		if err := r.GetTenantDBConnection(ctx).First(&departmentDomain.UserDepartment{}, "id = ?", req.GetParams().GetDepartment()).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Department not found")
		}
	}

	updates := make(map[string]interface{})
	if req.GetParams().GetDisplayName() != "" {
		updates["display_name"] = req.GetParams().GetDisplayName()
		updates["name"] = utils.ConvertToKebabCase(req.GetParams().GetDisplayName())
	}
	if req.GetParams().GetName() != "" {
		updates["name"] = req.GetParams().GetName()
	}
	if req.GetParams().GetDepartment() != 0 {
		updates["department_id"] = req.GetParams().GetDepartment()
	}

	if err := r.GetTenantDBConnection(ctx).Model(&externalRoleDomain.Role{}).Where("id = ?", req.GetParams().GetId()).Updates(updates).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	// Delete existing role permissions assignments and create new ones
	if !req.GetParams().GetEmptyPermissions() {
		if err := r.GetTenantDBConnection(ctx).Where("role_id = ?", req.GetParams().GetId()).Delete(&domain.RolePermissionAssignment{}).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		role := &externalRoleDomain.Role{}
		if err := r.GetTenantDBConnection(ctx).First(role, req.GetParams().GetId()).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		for _, permission := range req.GetParams().GetPermissions() {
			rolePermission := &domain.RolePermissionAssignment{
				RoleID:       role.ID,
				PermissionID: permission,
			}

			if err := r.GetTenantDBConnection(ctx).Create(rolePermission).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}

	}

	role := &externalRoleDomain.Role{}
	if err := r.GetTenantDBConnection(ctx).First(role, req.GetParams().GetId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var err error
	role, err = r.GetAllRoleData(ctx, role.ID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.UpdateRoleResponse{
		Result: convertors.BuildRoleResponse(role),
	}, nil
}
