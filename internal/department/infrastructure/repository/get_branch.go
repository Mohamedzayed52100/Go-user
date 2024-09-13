package repository

import (
	"context"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	"net/http"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"google.golang.org/grpc/status"
)

func (r *DepartmentRepository) GetLoggedInUser(ctx context.Context) (*externalUserDomain.User, error) {
	var (
		loggedInUser *externalUserDomain.User
		role         *externalRoleDomain.Role
	)

	authToken := ctx.Value(meta.AuthorizationContextKey.String())
	if authToken == nil {
		return nil, status.Error(http.StatusUnauthorized, "Unauthorized")
	}

	accMeta := auth.GetAccountMetadataFromToken(authToken.(string))

	if accMeta.Email == "" {
		return nil, nil
	}

	if err := r.GetSharedDB().
		Model(&externalUserDomain.User{}).
		Where("email = ?", accMeta.Email).
		First(&loggedInUser).Error; err != nil {
		return nil, err
	}

	role, err := r.GetAllRoleData(ctx, loggedInUser.RoleID)
	if err != nil {
		return nil, err
	}

	loggedInUser.Role = role

	return loggedInUser, nil
}

func (r *DepartmentRepository) GetCurrentBranchId(ctx context.Context) int32 {
	var branch int32

	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return 0
	}

	if currentUser == nil {
		authToken := ctx.Value(meta.AuthorizationContextKey.String())
		clientBranch := auth.GetClientBranchFromToken(authToken.(string))
		if clientBranch == 0 {
			r.GetTenantDBConnection(ctx).Model(externalUserDomain.Branch{}).Select("id").Scan(&branch)
			return branch
		}

		currentUser = &externalUserDomain.User{
			BranchID: clientBranch,
		}
	}

	r.GetTenantDBConnection(ctx).Model(&externalUserDomain.Branch{}).
		Where("id = ?", currentUser.BranchID).
		Select("id").
		Scan(&branch)

	return branch
}

func (r *DepartmentRepository) GetAllRoleData(ctx context.Context, id int32) (*externalRoleDomain.Role, error) {
	var role *externalRoleDomain.Role
	if err := r.GetTenantDBConnection(ctx).Where("id =?", id).First(&role).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var department *departmentDomain.UserDepartment
	if err := r.GetTenantDBConnection(ctx).Model(&departmentDomain.UserDepartment{}).Where("id = ?", role.DepartmentID).First(&department).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	role.Department = department

	permissions := []*roleDomain.Permission{}
	r.GetTenantDBConnection(ctx).Model(&roleDomain.Permission{}).Joins("JOIN role_permission_assignments ON role_permission_assignments.permission_id = permissions.id").Where("role_id =?", role.ID).Find(&permissions)
	role.Permissions = permissions

	// Count all users assigned to this role
	var count int64
	r.GetSharedDB().Model(&externalUserDomain.User{}).Where("role_id = ?", role.ID).Count(&count)
	role.UsersCount = int32(count)

	return role, nil
}
