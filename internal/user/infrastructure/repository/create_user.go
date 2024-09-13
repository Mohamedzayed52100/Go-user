package repository

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/goplaceapp/goplace-user/internal/user/adapters/converters"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"gorm.io/gorm"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/utils"
	"google.golang.org/grpc/status"
)

func (r *UserRepository) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.CreateUserResponse, error) {
	currentUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	joined, err := time.Parse(time.DateOnly, req.GetParams().GetJoinedAt())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	birthdate, err := time.Parse(time.DateOnly, req.GetParams().GetBirthdate())
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.SharedDbConnection.First(&externalUserDomain.User{}, "phone_number = ?", req.GetParams().GetPhoneNumber()).Error; err == nil {
		return nil, status.Error(http.StatusConflict, "This phone number is already in use")
	}

	if err := r.SharedDbConnection.First(&externalUserDomain.User{}, "email = ?", req.GetParams().GetEmail()).Error; err == nil {
		return nil, status.Error(http.StatusConflict, "This email address is already in use")
	}

	if err := r.SharedDbConnection.First(&externalUserDomain.User{}, "employee_id = ?", req.GetParams().GetEmployeeId()).Error; err == nil && req.GetParams().GetEmployeeId() != "" {
		return nil, status.Error(http.StatusConflict, "Employee ID already assigned")
	}

	generateUserPassword := utils.GeneratePassword()
	hashedPassword := utils.HashPassword(generateUserPassword)

	reqUser := &externalUserDomain.User{
		EmployeeID:  req.GetParams().GetEmployeeId(),
		FirstName:   req.GetParams().GetFirstName(),
		LastName:    req.GetParams().GetLastName(),
		Email:       strings.ToLower(req.GetParams().GetEmail()),
		Password:    hashedPassword,
		PhoneNumber: req.GetParams().GetPhoneNumber(),
		Avatar:      req.GetParams().GetAvatar(),
		RoleID:      req.GetParams().GetRole(),
		TenantID:    currentUser.TenantID,
		JoinedAt:    joined,
		Birthdate:   birthdate,
		PinCode:     utils.GeneratePinCode(),
		BranchID:    currentUser.BranchID,
	}

	if req.GetParams().GetBranchIds() == nil {
		req.Params.BranchIds = append(req.Params.BranchIds, currentUser.BranchID)
	}

	for _, v := range req.GetParams().GetBranchIds() {
		if !r.CheckForBranchAccess(ctx, v) {
			return nil, status.Error(http.StatusNotFound, "Please only assign branches to which you have access")
		}
	}

	if err := r.SharedDbConnection.Create(&reqUser).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	for _, branch := range req.GetParams().GetBranchIds() {
		if err := r.GetTenantDBConnection(ctx).Create(&externalUserDomain.UserBranchAssignment{
			UserID:   reqUser.ID,
			BranchID: branch,
		}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Error(http.StatusInternalServerError, "Branch not found")
			}
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}
	}

	if err = r.SendNewUserInvitationEmail(reqUser, generateUserPassword); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	reqUser, err = r.GetAllUserData(ctx, reqUser.ID)
	if err != nil {
		return nil, err
	}

	return &userProto.CreateUserResponse{
		Result: converters.BuildUserResponse(reqUser),
	}, nil
}
