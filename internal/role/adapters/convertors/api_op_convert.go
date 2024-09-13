package convertors

import (
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func BuildAllRolesResponse(roles []*externalRoleDomain.Role) []*userProto.Role {
	res := []*userProto.Role{}
	for _, role := range roles {
		res = append(res, BuildRoleResponse(role))
	}
	return res
}

func BuildPermissionResponse(permission *roleDomain.Permission) *userProto.Permission {
	res := &userProto.Permission{
		Id:          permission.ID,
		Name:        permission.Name,
		DisplayName: permission.DisplayName,
		Description: permission.Description,
		CreatedAt:   timestamppb.New(permission.CreatedAt),
		UpdatedAt:   timestamppb.New(permission.UpdatedAt),
	}

	return res
}

func BuildAllPermissionsResponse(permissions []*roleDomain.Permission) []*userProto.Permission {
	res := []*userProto.Permission{}
	for _, permission := range permissions {
		res = append(res, BuildPermissionResponse(permission))
	}

	return res
}
