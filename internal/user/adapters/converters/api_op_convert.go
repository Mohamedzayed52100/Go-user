package converters

import (
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildAuthenticatedUserResponse(user *externalUserDomain.User) *userProto.AuthenticatedUser {
	return &userProto.AuthenticatedUser{
		Id:          user.ID,
		EmployeeId:  user.EmployeeID,
		BirthDate:   timestamppb.New(user.Birthdate),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        BuildRoleResponse(user.Role),
		Department:  BuildDepartmentResponse(user.Role.Department),
		Gender:      user.Gender,
		Timezone:    user.Timezone,
		GroupId:     user.TenantID,
		Branch: &userProto.UBranch{
			Id:       int32(user.Branch.ID),
			Name:     user.Branch.Name,
			Currency: user.Branch.Currency,
		},
		Avatar: user.Avatar,
	}
}

func BuildUserResponse(user *externalUserDomain.User) *userProto.User {
	res := &userProto.User{
		Id:          user.ID,
		EmployeeId:  user.EmployeeID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.DisplayName,
		Gender:      user.Gender,
		Department:  user.Role.Department.Name,
		Avatar:      user.Avatar,
		JoinedAt:    timestamppb.New(user.JoinedAt),
		Birthdate:   timestamppb.New(user.Birthdate),
	}

	if user.Branches != nil {
		for _, branch := range user.Branches {
			res.Branches = append(res.Branches, &userProto.UserBranch{
				Id:   int32(branch.ID),
				Name: branch.Name,
			})
		}
	}

	return res
}

func BuildAllUsersResponse(users []*externalUserDomain.User) []*userProto.User {
	res := []*userProto.User{}
	for _, user := range users {
		res = append(res, BuildUserResponse(user))
	}
	return res
}

func BuildRoleResponse(role *externalRoleDomain.Role) *userProto.Role {
	res := &userProto.Role{
		Id:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Department:  role.Department.Name,
		UsersCount:  role.UsersCount,
		CreatedAt:   timestamppb.New(role.CreatedAt),
		UpdatedAt:   timestamppb.New(role.UpdatedAt),
	}

	if role.Permissions != nil {
		for _, permission := range role.Permissions {
			res.Permissions = append(res.Permissions, permission.Name)
		}
	}

	return res
}

func BuildDepartmentResponse(department *departmentDomain.UserDepartment) *userProto.Department {
	return &userProto.Department{
		Id:        department.ID,
		Name:      department.Name,
		CreatedAt: timestamppb.New(department.CreatedAt),
		UpdatedAt: timestamppb.New(department.UpdatedAt),
	}
}

func BuildAllBranchesResponse(branches []*externalUserDomain.Branch) []*userProto.Branch {
	result := make([]*userProto.Branch, 0)
	for _, branch := range branches {
		result = append(result, BuildBranchResponse(branch))

	}

	return result
}

func BuildBranchResponse(branch *externalUserDomain.Branch) *userProto.Branch {
	return &userProto.Branch{
		Id:          int32(branch.ID),
		Name:        branch.Name,
		Country:     branch.Country,
		City:        branch.City,
		Address:     branch.Address,
		Email:       branch.Email,
		PhoneNumber: branch.PhoneNumber,
		Website:     branch.Website,
	}
}
