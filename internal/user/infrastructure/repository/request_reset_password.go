package repository

import (
	"context"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"

	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) RequestResetPassword(ctx context.Context, req *userProto.RequestResetPasswordRequest) (*userProto.RequestResetPasswordResponse, error) {
	var (
		currentUser  *domain.User
		verifyMethod string
	)

	if err := r.GetSharedDB().
		First(&currentUser, "email = ? OR phone_number = ?", req.GetEmail(), req.GetPhoneNumber()).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, err.Error())
	}

	if req.GetEmail() != "" {
		verifyMethod = "email"
	} else if req.GetPhoneNumber() != "" {
		verifyMethod = "phone_number"
	} else {
		return nil, status.Error(http.StatusBadRequest, "Email or phone number is required")
	}

	otp, token, err := r.SendUserOTP(
		ctx,
		currentUser,
		verifyMethod,
		"reset_password",
		"",
		4,
	)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if req.GetEmail() != "" {
		_, err := r.SendResetPasswordEmail(currentUser, token, otp)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Failed to send otp via email, please try again")
		}
	} else if req.GetPhoneNumber() != "" {
		if err := r.SendResetPasswordWhatsappMessage(ctx, currentUser, otp); err != nil {
			return nil, status.Error(http.StatusInternalServerError, "Failed to send otp via whatsapp, please try again")
		}
	}

	return &userProto.RequestResetPasswordResponse{
		Token: token,
	}, nil
}
