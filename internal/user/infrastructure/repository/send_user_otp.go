package repository

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/goplaceapp/goplace-user/utils"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) SendUserOTP(ctx context.Context, user *domain.User, verifyMethod, otpType, token string, length int32) (string, string, error) {
	if length == 0 {
		length = 4
	}

	var otp string
	for i := 0; i < int(length); i++ {
		otp += strconv.FormatInt(int64(rand.Intn(10)), 10)
	}

	if token == "" {
		token = utils.GenerateRandomString(32)
	}

	if err := r.GetSharedDB().Create(&domain.UserOtp{
		UserID:       user.ID,
		Code:         otp,
		Token:        token,
		Type:         otpType,
		VerifyMethod: verifyMethod,
		ExpiresAt:    time.Now().UTC().Add(15 * time.Minute),
	}).Error; err != nil {
		if err := r.GetSharedDB().
			Where("user_id = ?", user.ID).
			Updates(&domain.UserOtp{
				Code:      otp,
				Token:     token,
				ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
			}).Error; err != nil {
			return "", "", status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return otp, token, nil
}
