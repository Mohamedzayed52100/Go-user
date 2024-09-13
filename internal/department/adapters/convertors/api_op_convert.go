package convertors

import (
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BuildDepartmentResponse(department *domain.UserDepartment) *userProto.Department {
	return &userProto.Department{
		Id:        department.ID,
		Name:      department.Name,
		CreatedAt: timestamppb.New(department.CreatedAt),
		UpdatedAt: timestamppb.New(department.UpdatedAt),
	}
}

func BuildAllDepartmentsResponse(departments []*domain.UserDepartment) []*userProto.Department {
	res := []*userProto.Department{}
	for _, department := range departments {
		res = append(res, BuildDepartmentResponse(department))
	}
	return res
}
