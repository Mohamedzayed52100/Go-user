package grpc

import (
	"context"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

func (s *DepartmentServiceServer) GetAllDepartments(ctx context.Context, req *userProto.GetAllDepartmentsRequest) (*userProto.GetAllDepartmentsResponse, error) {
	res, err := s.departmentService.GetAllDepartments(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
