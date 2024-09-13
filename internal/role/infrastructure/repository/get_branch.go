package repository

import (
	"context"
	"net/http"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"google.golang.org/grpc/status"
)

func (r *RoleRepository) GetAllUserBranchesIDs(ctx context.Context) []int32 {
	branches := []int32{}

	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil
	}

	if currentUser == nil {
		authToken := ctx.Value(meta.AuthorizationContextKey.String())
		clientBranch := auth.GetClientBranchFromToken(authToken.(string))
		if clientBranch == 0 {
			r.GetTenantDBConnection(ctx).Model(externalUserDomain.Branch{}).Pluck("id", &branches)
		} else {
			branches = append(branches, clientBranch)
		}

		return branches
	}

	r.GetTenantDBConnection(ctx).
		Model(externalUserDomain.UserBranchAssignment{}).
		Where("user_id = ?", currentUser.ID).
		Distinct().
		Pluck("branch_id", &branches)

	return branches
}

func (r *RoleRepository) GetLoggedInUser(ctx context.Context) (*externalUserDomain.User, error) {
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

func (r *RoleRepository) GetCurrentBranchId(ctx context.Context) int32 {
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
