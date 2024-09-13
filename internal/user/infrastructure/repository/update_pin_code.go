package repository

import (
	"context"
	"net/http"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) UpdatePinCode(ctx context.Context, req *userProto.UpdatePinCodeRequest) (*userProto.UpdatePinCodeResponse, error) {
	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if currentUser.PinCode != req.GetOldPinCode() {
		return nil, status.Error(http.StatusBadRequest, "Old pin code is incorrect")
	}

	if req.GetConfirmNewPinCode() != req.GetNewPinCode() {
		return nil, status.Error(http.StatusBadRequest, "New pin code and confirmation new pin code does not match!")
	}

	if err := r.SharedDbConnection.Model(currentUser).Update("pin_code", req.GetNewPinCode()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.UpdatePinCodeResponse{
		Code:    http.StatusOK,
		Message: "Pin code changed successfully!",
	}, nil
}
