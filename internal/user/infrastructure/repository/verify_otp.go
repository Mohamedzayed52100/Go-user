package repository

import (
	"context"
	"net/http"
	"time"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) VerifyOtp(ctx context.Context, req *userProto.VerifyOtpRequest) (*userProto.VerifyOtpResponse, error) {
	var otp domain.UserOtp
	if err := r.GetSharedDB().First(&otp, "token = ? and code = ?", req.GetToken(), req.GetCode()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Invalid OTP code")
	}

	if otp.ExpiresAt.Before(time.Now().UTC()) {
		return nil, status.Error(http.StatusConflict, "Your code has expired! Please request a new one.")
	}

	if err := r.GetSharedDB().Table("user_otps").Where("token = ? and code = ?", req.GetToken(), req.GetCode()).Update("verified", "true").Error; err != nil {
		return nil, status.Error(http.StatusConflict, err.Error())
	}

	return &userProto.VerifyOtpResponse{
		Code:    http.StatusOK,
		Message: "OTP code is verified!",
	}, nil
}
