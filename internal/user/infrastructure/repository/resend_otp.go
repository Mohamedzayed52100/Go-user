package repository

import (
	"context"
	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) ResendOtp(ctx context.Context, req *userProto.ResendOtpRequest) (*userProto.ResendOtpResponse, error) {
	var (
		newOTP  string
		userOTP *domain.UserOtp
		user    *domain.User
	)

	r.GetSharedDB().Table("user_otps").
		Where("token = ?", req.GetToken()).
		First(&userOTP)

	for i := 0; i < len(userOTP.Code); i++ {
		newOTP += strconv.FormatInt(int64(rand.Intn(10)), 10)
	}

	if err := r.GetSharedDB().Table("user_otps").
		Where("token = ?", req.GetToken()).
		Updates(map[string]interface{}{
			"code":       newOTP,
			"expires_at": time.Now().UTC().Add(15 * time.Minute),
		}).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetSharedDB().Model(&domain.User{}).
		Where("id = ?", userOTP.UserID).
		First(&user).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	switch userOTP.Type {
	case "reset_password":
		if userOTP.VerifyMethod == "email" {
			_, err := r.SendResetPasswordEmail(user, userOTP.Token, newOTP)
			if err != nil {
				return nil, status.Error(http.StatusInternalServerError, "Failed to send otp via email, please try again")
			}
		} else if userOTP.VerifyMethod == "phone_number" {
			// TODO: Send whatsapp notification
		}
	default:
		return nil, status.Error(http.StatusBadRequest, "Invalid OTP type")
	}

	return &userProto.ResendOtpResponse{
		Code:    http.StatusOK,
		Message: "OTP sent successfully",
	}, nil
}
