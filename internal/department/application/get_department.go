package application

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *DepartmentService) GetAllDepartments(ctx context.Context, req *userProto.GetAllDepartmentsRequest) (*userProto.GetAllDepartmentsResponse, error) {
	res, err := s.Repository.GetAllDepartments(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
