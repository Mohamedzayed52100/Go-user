package repository

import (
	"context"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
	"net/http"
)

func (r *UserRepository) ResetPassword(ctx context.Context, req *userProto.ResetPasswordRequest) (*userProto.ResetPasswordResponse, error) {
	var userId int32

	if err := r.GetSharedDB().
		Table("user_otps").
		Where(`token = ? AND verified = 'true'`, req.GetToken()).
		Select("user_id").
		Scan(&userId).Error; err != nil || userId == 0 {
		return nil, status.Error(http.StatusNotFound, "Invalid token provided!")
	}

	if req.GetConfirmNewPassword() != req.GetNewPassword() {
		return nil, status.Error(http.StatusBadRequest, "New password and confirmation does not match!")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.SharedDbConnection.
		Table("users").
		Where("id = ?", userId).
		Update("password", password).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.ResetPasswordResponse{
		Code:    http.StatusOK,
		Message: "Password changed successfully!",
	}, nil
}
