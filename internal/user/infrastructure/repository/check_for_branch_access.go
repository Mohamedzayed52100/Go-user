package repository

import (
	"context"

	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
)

func (r *UserRepository) CheckForBranchAccess(ctx context.Context, branchId int32) bool {
	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return false
	}

	if currentUser == nil {
		authToken := ctx.Value(meta.AuthorizationContextKey.String())

		clientBranch := auth.GetClientBranchFromToken(authToken.(string))
		if clientBranch == 0 ||
			clientBranch == branchId {
			return true
		}

		return clientBranch == branchId
	}

	if err := r.GetTenantDBConnection(ctx).First(&externalUserDomain.UserBranchAssignment{}, "user_id = ? AND branch_id = ?", currentUser.ID, branchId).Error; err != nil {
		return false
	}

	return true
}
