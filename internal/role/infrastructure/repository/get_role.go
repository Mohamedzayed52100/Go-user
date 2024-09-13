package repository

import (
	"context"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
	"net/http"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/common"
	"github.com/goplaceapp/goplace-user/internal/role/adapters/convertors"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *RoleRepository) GetAllRoles(ctx context.Context, req *userProto.GetAllRolesRequest) (*userProto.GetAllRolesResponse, error) {
	var (
		count int64
		roles = []*externalRoleDomain.Role{}
	)

	offset := (req.GetParams().GetCurrentPage() - 1) * req.GetParams().GetPerPage()
	limit := req.GetParams().GetPerPage()

	r.RolesQueryBuilder(ctx, req).Model(&externalRoleDomain.Role{}).Count(&count)

	r.RolesQueryBuilder(ctx, req).
		Offset(int(offset)).
		Limit(int(limit)).
		Order("id DESC").
		Find(&roles)

	for i := range roles {
		var err error
		roles[i], err = r.GetAllRoleData(ctx, roles[i].ID)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &userProto.GetAllRolesResponse{
		Pagination: common.BuildPaginationResponse(req.GetParams(), limit, offset),
		Result:     convertors.BuildAllRolesResponse(roles),
	}, nil
}

func (r *RoleRepository) RolesQueryBuilder(ctx context.Context, req *userProto.GetAllRolesRequest) *gorm.DB {
	query := r.GetTenantDBConnection(ctx).Model(&externalRoleDomain.Role{})

	if len(req.GetDepartment()) > 0 {
		departments := []int32{}

		r.GetTenantDBConnection(ctx).
			Model(&departmentDomain.UserDepartment{}).
			Where("id IN ?", req.GetDepartment()).
			Distinct().Pluck("id", &departments)

		query = query.Where("department_id IN (?)", departments)
	}

	if req.GetQuery() != "" {
		searchQuery := "%" + req.GetQuery() + "%"
		query = query.Where("(UPPER(name) LIKE UPPER(?) OR UPPER(display_name) LIKE UPPER(?))", searchQuery, searchQuery)
	}
	return query
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, req *userProto.GetRoleByIDRequest) (*userProto.GetRoleByIDResponse, error) {
	var role *externalRoleDomain.Role
	var err error
	if err = r.GetTenantDBConnection(ctx).First(&role, req.GetId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err = r.GetTenantDBConnection(ctx).First(&role, "id = ?", req.GetId()).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "Role not found")
	}

	role, err = r.GetAllRoleData(ctx, role.ID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.GetRoleByIDResponse{
		Result: convertors.BuildRoleResponse(role),
	}, nil
}
