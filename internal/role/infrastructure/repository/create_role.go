package repository

import (
	"context"
	"net/http"

	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/role/adapters/convertors"
	"github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"github.com/goplaceapp/goplace-user/utils"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *RoleRepository) CreateRole(ctx context.Context, req *userProto.CreateRoleRequest) (*userProto.CreateRoleResponse, error) {
	name := utils.ConvertToKebabCase(req.GetParams().GetDisplayName())

	if err := r.GetTenantDBConnection(ctx).First(&departmentDomain.UserDepartment{}, "id = ?", req.GetParams().GetDepartment()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Branch not found")
	}

	if err := r.GetTenantDBConnection(ctx).First(&departmentDomain.UserDepartment{}, "id = ?", req.GetParams().GetDepartment()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "You don't have access to this branch")
	}

	role := &externalRoleDomain.Role{
		DepartmentID: req.GetParams().GetDepartment(),
		Name:         name,
		DisplayName:  req.GetParams().GetDisplayName(),
	}

	if err := r.GetTenantDBConnection(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return status.Error(http.StatusInternalServerError, err.Error())
		}

		for _, permission := range req.GetParams().GetPermissions() {
			rolePermission := &domain.RolePermissionAssignment{
				RoleID:       role.ID,
				PermissionID: permission,
			}

			if err := tx.Create(rolePermission).Error; err != nil {
				return status.Error(http.StatusInternalServerError, err.Error())
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	var err error
	role, err = r.GetAllRoleData(ctx, role.ID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.CreateRoleResponse{
		Result: convertors.BuildRoleResponse(role),
	}, nil
}
