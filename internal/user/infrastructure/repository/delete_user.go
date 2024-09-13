package repository

import (
	"context"
	"net/http"

	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) DeleteUser(ctx context.Context, req *userProto.DeleteUserRequest) (*userProto.DeleteUserResponse, error) {
	currentBranchID := r.GetCurrentBranchId(ctx)

	if err := r.GetTenantDBConnection(ctx).
		Model(externalUserDomain.UserBranchAssignment{}).
		Where("user_id = ? AND branch_id = ?", req.GetId(), currentBranchID).
		First(&externalUserDomain.UserBranchAssignment{}).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "User not found")
	}

	if err := r.SharedDbConnection.Delete(&externalUserDomain.User{}, "id = ?", req.GetId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.DeleteUserResponse{
		Code:    http.StatusOK,
		Message: "Deleted user successfully",
	}, nil
}
