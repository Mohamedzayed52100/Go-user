package repository

import (
	"context"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) ChangePassword(ctx context.Context, req *userProto.ChangePasswordRequest) (*userProto.ChangePasswordResponse, error) {
	var isSetupCompleted bool

	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.GetOldPassword())); err != nil {
		return nil, status.Error(http.StatusUnauthorized, "Old password does not match!")
	}

	if req.GetConfirmNewPassword() != req.GetNewPassword() {
		return nil, status.Error(http.StatusBadRequest, "New password and confirmation new password does not match!")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	currentUser.Password = string(password)

	if err := r.SharedDbConnection.Model(&currentUser).Update("password", password).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if currentUser.SetupCompleted == false {
		isSetupCompleted = true

		if err := r.SharedDbConnection.
			Model(&currentUser).
			Update("setup_completed", isSetupCompleted).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &userProto.ChangePasswordResponse{
		Code:    http.StatusOK,
		Message: "Password changed successfully!",
	}, nil
}
