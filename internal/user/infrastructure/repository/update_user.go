package repository

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/goplaceapp/goplace-user/internal/user/adapters/converters"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *UserRepository) UpdateUser(ctx context.Context, req *userProto.UpdateUserRequest) (*userProto.UpdateUserResponse, error) {
	var u externalUserDomain.User
	currentBranchID := r.GetCurrentBranchId(ctx)

	if err := r.GetTenantDBConnection(ctx).
		Model(externalUserDomain.UserBranchAssignment{}).
		Where("user_id = ? AND branch_id = ?", req.GetParams().GetId(), currentBranchID).
		First(&externalUserDomain.UserBranchAssignment{}).
		Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "User not found")
	}

	if err := r.SharedDbConnection.First(&u, "id = ?", req.GetParams().GetId()).Error; err != nil {
		return nil, status.Error(http.StatusNotFound, "User not found")
	}

	for _, v := range req.GetParams().GetBranchIds() {
		if !r.CheckForBranchAccess(ctx, v) {
			return nil, status.Error(http.StatusNotFound, "Please only assign branches to which you have access")
		}
	}

	if req.GetParams().GetPhoneNumber() != "" && req.GetParams().GetPhoneNumber() != u.PhoneNumber {
		if err := r.SharedDbConnection.First(&externalUserDomain.User{}, "phone_number = ?", req.GetParams().GetPhoneNumber()).Error; err == nil {
			return nil, status.Error(http.StatusConflict, "This phone number is already in use")
		}
	}

	if req.GetParams().GetEmail() != "" && req.GetParams().GetEmail() != u.Email {
		if err := r.SharedDbConnection.
			First(&externalUserDomain.User{}, "email = ?", strings.ToLower(req.GetParams().GetEmail())).
			Error; err == nil {
			return nil, status.Error(http.StatusConflict, "This email address is already in use")
		}
	}

	if req.GetParams().GetEmployeeId() != "" && req.GetParams().GetEmployeeId() != u.EmployeeID {
		if err := r.SharedDbConnection.First(&externalUserDomain.User{}, "employee_id = ?", req.GetParams().GetEmployeeId()).Error; err == nil && req.GetParams().GetEmployeeId() != "" {
			return nil, status.Error(http.StatusConflict, "Employee ID already assigned")
		}
	}

	updates := make(map[string]interface{})
	if req.GetParams().GetAvatar() != "" {
		updates["avatar"] = req.GetParams().GetAvatar()
	}
	if req.GetParams().GetFirstName() != "" {
		updates["first_name"] = req.GetParams().GetFirstName()
	}
	if req.GetParams().GetLastName() != "" {
		updates["last_name"] = req.GetParams().GetLastName()
	}
	if req.GetParams().GetPassword() != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(req.GetParams().GetPassword()), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		updates["password"] = password
	}
	if req.GetParams().GetRole() != 0 {
		updates["role_id"] = req.GetParams().GetRole()
	}
	if req.GetParams().GetBirthdate() != "" {
		birthdate, err := time.Parse(time.DateOnly, req.GetParams().GetBirthdate())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		updates["birthdate"] = birthdate
	}
	if req.GetParams().GetJoinedAt() != "" {
		joinedAt, err := time.Parse(time.DateOnly, req.GetParams().GetJoinedAt())
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
		updates["joined_at"] = joinedAt
	}
	if req.GetParams().GetEmployeeId() != "" {
		updates["employee_id"] = req.GetParams().GetEmployeeId()
	}

	if err := r.SharedDbConnection.Model(&externalUserDomain.User{}).Where("id = ?", req.GetParams().GetId()).Updates(updates).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if req.GetParams().GetEmail() != "" && req.GetParams().GetEmail() != u.Email {
		if err := r.SharedDbConnection.
			Model(&externalUserDomain.User{}).
			Where("id =?", req.GetParams().GetId()).
			Update("email", strings.ToLower(req.GetParams().GetEmail())).
			Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	if req.GetParams().GetPhoneNumber() != "" && req.GetParams().GetPhoneNumber() != u.PhoneNumber {
		if err := r.SharedDbConnection.Model(&externalUserDomain.User{}).Where("id =?", req.GetParams().GetId()).Update("phone_number", req.GetParams().GetPhoneNumber()).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, status.Error(http.StatusInternalServerError, "Phone number already exists")
			}
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	if req.GetParams().GetBranchIds() != nil {
		existingBranches := []int32{}
		r.GetTenantDBConnection(ctx).
			Model(externalUserDomain.UserBranchAssignment{}).
			Where("user_id = ?", req.GetParams().GetId()).
			Pluck("branch_id", &existingBranches)

		existingBranchesMap := make(map[int32]bool)
		for _, branchId := range existingBranches {
			existingBranchesMap[branchId] = true
		}

		for _, branchId := range req.GetParams().GetBranchIds() {
			if _, ok := existingBranchesMap[branchId]; !ok {
				if err := r.GetTenantDBConnection(ctx).Model(externalUserDomain.UserBranchAssignment{}).Create(&externalUserDomain.UserBranchAssignment{
					UserID:   req.GetParams().GetId(),
					BranchID: branchId,
				}).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return nil, status.Error(http.StatusInternalServerError, "Branch not found")
					}
					if !errors.Is(err, gorm.ErrDuplicatedKey) {
						return nil, status.Error(http.StatusInternalServerError, err.Error())
					}
				}
			}
		}

		newBranchesMap := make(map[int32]bool)
		for _, branch := range req.GetParams().GetBranchIds() {
			newBranchesMap[branch] = true
		}

		for _, t := range existingBranches {
			if _, ok := newBranchesMap[t]; !ok {
				if err := r.GetTenantDBConnection(ctx).
					Model(externalUserDomain.UserBranchAssignment{}).
					Delete(&externalUserDomain.UserBranchAssignment{}, "user_id = ? AND branch_id = ?",
						req.GetParams().GetId(), t,
					).
					Error; err != nil {
					return nil, status.Error(http.StatusInternalServerError, err.Error())
				}
			}
		}
	}

	updatedUser, err := r.GetAllUserData(ctx, req.GetParams().GetId())
	if err != nil {
		return nil, err
	}

	if len(req.GetParams().GetBranchIds()) > 0 {
		if err := r.GetSharedDB().Table("users").Where("id = ?", req.GetParams().GetId()).Update("branch_id = ?", req.Params.BranchIds[0]).Error; err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &userProto.UpdateUserResponse{
		Result: converters.BuildUserResponse(updatedUser),
	}, nil
}
