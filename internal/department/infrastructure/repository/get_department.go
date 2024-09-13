package repository

import (
	"context"
	"github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/common"
	"github.com/goplaceapp/goplace-user/internal/department/adapters/convertors"
)

func (r *DepartmentRepository) GetAllDepartments(ctx context.Context, req *userProto.GetAllDepartmentsRequest) (*userProto.GetAllDepartmentsResponse, error) {
	var (
		count       int64
		departments []*domain.UserDepartment
	)

	offset := (req.GetParams().GetCurrentPage() - 1) * req.GetParams().GetPerPage()
	limit := req.GetParams().GetPerPage()

	if req.GetParams().GetPerPage() <= 0 || req.GetParams().GetCurrentPage() <= 0 {
		limit = 1e9
		offset = 0
	}

	r.GetTenantDBConnection(ctx).Model(&domain.UserDepartment{}).Count(&count)
	r.GetTenantDBConnection(ctx).Offset(int(offset)).Limit(int(limit)).Order("id DESC").Find(&departments)

	return &userProto.GetAllDepartmentsResponse{
		Pagination: common.BuildPaginationResponse(req.GetParams(), int32(count), offset),
		Result:     convertors.BuildAllDepartmentsResponse(departments),
	}, nil
}
