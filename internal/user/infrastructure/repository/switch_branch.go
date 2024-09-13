package repository

import (
	"context"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) SwitchBranch(ctx context.Context, req *userProto.SwitchBranchRequest) (*userProto.SwitchBranchResponse, error) {
	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).First(&domain.Branch{}, "id = ?", req.GetBranchId()).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "Branch not found")
	}

	if !r.CheckForBranchAccess(ctx, req.GetBranchId()) {
		return nil, status.Error(http.StatusNotFound, "You don't have access to this branch")
	}

	currentUser.BranchID = req.GetBranchId()
	if err := r.SharedDbConnection.Table("users").Where("id = ?", currentUser.ID).Update("branch_id", currentUser.BranchID).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.SwitchBranchResponse{
		Code:    http.StatusOK,
		Message: "Branch switched successfully",
	}, nil
}
