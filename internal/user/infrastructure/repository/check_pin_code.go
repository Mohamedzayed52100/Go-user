package repository

import (
	"context"
	"net/http"

	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) CheckPinCode(ctx context.Context, req *userProto.PinCodeRequest) (*userProto.PinCodeResponse, error) {
	var (
		u       externalUserDomain.User
		pinCode = req.GetPinCode()
	)

	if req.GetRole() == "user" || req.GetRole() == "" {
		email := auth.GetAccountMetadataFromToken(ctx.Value(meta.AuthorizationContextKey.String()).(string)).Email
		if email == "" {
			return nil, status.Error(http.StatusNotFound, errorhelper.ErrEmailAddressNotFound)
		}

		if err := r.SharedDbConnection.
			Where("email = ?", email).
			First(&u).Error; err != nil {
			return nil, status.Error(http.StatusNotFound, err.Error())
		}

		if u.PinCode != pinCode {
			return nil, status.New(http.StatusForbidden, errorhelper.ErrIncorrectPinCode).Err()
		}
	} else if req.GetRole() == "su" {
		if !r.CheckAdminPinCode(ctx, req.GetPinCode()) {
			return nil, status.New(http.StatusForbidden, errorhelper.ErrIncorrectPinCode).Err()
		}
	}

	return &userProto.PinCodeResponse{
		Result: &userProto.PinCodeResult{
			Status:  "OK",
			Message: "Successfully checked pin code",
		},
	}, nil
}
