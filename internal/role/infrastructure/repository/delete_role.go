package repository

import (
	"context"
	"net/http"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/grpc/status"
)

func (r *RoleRepository) DeleteRole(ctx context.Context, req *userProto.DeleteRoleRequest) (*userProto.DeleteRoleResponse, error) {
	if err := r.GetTenantDBConnection(ctx).
		First(
			&externalRoleDomain.Role{}, "id = ?",
			req.GetId(),
		).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Role not found")
	}

	if err := r.GetTenantDBConnection(ctx).Delete(&externalRoleDomain.Role{}, "id = ?", req.GetId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Role is used by users")
	}

	return &userProto.DeleteRoleResponse{
		Code:    http.StatusOK,
		Message: "Role deleted successfully",
	}, nil
}
