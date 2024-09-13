package repository

import (
	"context"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/user/adapters/converters"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
)

func (r *UserRepository) GetCurrentBranchId(ctx context.Context) int32 {
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

func (r *UserRepository) GetBranchByID(ctx context.Context, branchID int32) (*externalUserDomain.Branch, error) {
	var branch *externalUserDomain.Branch

	err := r.GetTenantDBConnection(ctx).Where("id =?", branchID).First(&branch).Error
	if err != nil {
		return branch, err
	}

	return branch, nil
}

func (r *UserRepository) GetAllUserBranchesIDs(ctx context.Context) []int32 {
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

func (r *UserRepository) GetAllBranches(ctx context.Context, req *emptypb.Empty) (*userProto.GetAllBranchesResponse, error) {
	var branches []*externalUserDomain.Branch

	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		currentUser = nil
	}

	if currentUser == nil {
		if err := r.GetTenantDBConnection(ctx).Find(&branches).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	} else if currentUser != nil {
		if err := r.GetTenantDBConnection(ctx).
			Model(&externalUserDomain.Branch{}).
			Joins("JOIN user_branch_assignments ON user_branch_assignments.branch_id = branches.id").
			Find(&branches, "user_branch_assignments.user_id = ?", currentUser.ID).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &userProto.GetAllBranchesResponse{
		Result: converters.BuildAllBranchesResponse(branches),
	}, nil
}
