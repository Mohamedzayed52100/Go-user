package repository

import (
	"context"
	"net/http"

	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"

	"github.com/goplaceapp/goplace-user/internal/user/adapters/converters"
	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/common"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (r *UserRepository) GetAllUsers(ctx context.Context, req *userProto.GetAllUsersRequest) (*userProto.GetAllUsersResponse, error) {
	offset := (req.GetParams().GetCurrentPage() - 1) * req.GetParams().GetPerPage()
	limit := req.GetParams().GetPerPage()
	var count int64

	usersQueryBuilder(r, ctx, req).
		Model(&externalUserDomain.User{}).
		Count(&count)

	users := []*externalUserDomain.User{}
	usersQueryBuilder(r, ctx, req).
		Offset(int(offset)).
		Limit(int(limit)).
		Order("id DESC").
		Find(&users)

	for i := range users {
		var err error

		users[i], err = r.GetAllUserData(ctx, users[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return &userProto.GetAllUsersResponse{
		Pagination: common.BuildPaginationResponse(req.GetParams(), int32(count), offset),
		Result:     converters.BuildAllUsersResponse(users),
	}, nil
}

func usersQueryBuilder(r *UserRepository, ctx context.Context, req *userProto.GetAllUsersRequest) *gorm.DB {
	var (
		departmentIds []int
		roleIds       []int
	)

	usersDb := r.SharedDbConnection
	tenantDb := r.GetTenantDBConnection(ctx)

	query := usersDb.Model(&externalUserDomain.User{})

	userIds := []int{}
	r.GetTenantDBConnection(ctx).
		Model(externalUserDomain.UserBranchAssignment{}).
		Where("branch_id = ?", r.GetCurrentBranchId(ctx)).
		Pluck("user_id", &userIds)

	query = query.Where("id IN ?", userIds)

	if len(req.GetDepartment()) > 0 {
		tenantDb.Model(&departmentDomain.UserDepartment{}).Where("id IN (?)", req.GetDepartment()).Pluck("id", &departmentIds)
	}

	if len(departmentIds) > 0 {
		tenantDb.Model(&externalRoleDomain.Role{}).Where("department_id IN (?)", departmentIds).Pluck("id", &roleIds)
	}

	if len(req.GetRole()) > 0 {
		tenantDb.Model(&externalRoleDomain.Role{}).Where("id IN (?)", req.GetRole()).Pluck("id", &roleIds)
	}

	if len(roleIds) > 0 {
		query = query.Where("role_id IN (?)", roleIds)
	}

	if req.GetFromDate() != "" {
		query = query.Where("TO_CHAR(joined_at, 'YYYY-MM-DD') >= ?", req.GetFromDate())
	}

	if req.GetToDate() != "" {
		query = query.Where("TO_CHAR(joined_at, 'YYYY-MM-DD') <=?", req.GetToDate())
	}

	if req.GetQuery() != "" {
		searchQuery := "%" + req.GetQuery() + "%"
		query = query.Where("(UPPER(first_name) LIKE UPPER(?) OR "+
			"UPPER(last_name) LIKE UPPER(?) OR "+
			"phone_number LIKE ? OR "+
			"UPPER(email) LIKE UPPER(?))",
			searchQuery, searchQuery, searchQuery, searchQuery)
	}

	return query
}

func (r *UserRepository) GetUserByID(ctx context.Context, req *userProto.GetUserByIDRequest) (*userProto.GetUserByIDResponse, error) {
	user, err := r.GetAllUserData(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(http.StatusNotFound, "User not found")
	}

	for _, b := range user.Branches {
		if int32(b.ID) == r.GetCurrentBranchId(ctx) {
			return &userProto.GetUserByIDResponse{
				Result: converters.BuildUserResponse(user),
			}, nil
		}
	}

	return nil, status.Error(http.StatusNotFound, "User not found")
}

func (r *UserRepository) GetAuthenticatedUser(ctx context.Context) (*userProto.GetAuthenticatedUserResponse, error) {
	var (
		u       externalUserDomain.User
		country externalUserDomain.Country
	)

	getLoggedInUser, err := r.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	u = *getLoggedInUser

	getBranch, err := r.GetBranchByID(ctx, u.BranchID)
	if err != nil {
		return nil, err
	}

	if err := r.GetSharedDB().
		Model(&externalUserDomain.Country{}).
		Where("country_name = ?", getBranch.Country).
		First(&country).Error; err != nil {
		return nil, err
	}

	u.Timezone = country.Timezone
	u.Branch = &externalUserDomain.Branch{
		ID:        getBranch.ID,
		Name:      getBranch.Name,
		Currency: country.Currency,
		CreatedAt: getBranch.CreatedAt,
	}

	return &userProto.GetAuthenticatedUserResponse{
		Result: converters.BuildAuthenticatedUserResponse(&u),
	}, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*externalUserDomain.User, error) {
	var u externalUserDomain.User

	if err := r.SharedDbConnection.Where("email", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetLoggedInUser(ctx context.Context) (*externalUserDomain.User, error) {
	var (
		loggedInUser *externalUserDomain.User
		role         *externalRoleDomain.Role
	)

	authToken := ctx.Value(meta.AuthorizationContextKey.String())
	if authToken == nil {
		return nil, status.Error(http.StatusUnauthorized, "Unauthorized")
	}

	accMeta := auth.GetAccountMetadataFromToken(authToken.(string))

	if accMeta.Email == "" {
		return nil, nil
	}

	if err := r.GetSharedDB().
		Model(&externalUserDomain.User{}).
		Where("email = ?", accMeta.Email).
		First(&loggedInUser).Error; err != nil {
		return nil, err
	}

	role, err := r.RoleRepository.GetAllRoleData(ctx, loggedInUser.RoleID)
	if err != nil {
		return nil, err
	}

	loggedInUser.Role = role

	return loggedInUser, nil
}

func (r *UserRepository) GetAllUserData(ctx context.Context, id int32) (*externalUserDomain.User, error) {
	var reqUser *externalUserDomain.User
	if err := r.SharedDbConnection.Where("id = ?", id).First(&reqUser).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var role *externalRoleDomain.Role
	if err := r.GetTenantDBConnection(ctx).Where("id = ?", reqUser.RoleID).First(&role).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	reqUser.Role = role

	var department *departmentDomain.UserDepartment
	if err := r.GetTenantDBConnection(ctx).Where("id = ?", role.DepartmentID).First(&department).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}
	reqUser.Role.Department = department

	branches := []*externalUserDomain.Branch{}
	if err := r.GetTenantDBConnection(ctx).
		Model(branches).
		Joins("JOIN user_branch_assignments ON user_branch_assignments.branch_id = branches.id").
		Where("user_id = ?", id).
		Select("branches.*").
		Find(&branches).Error; err != nil {
		return nil, status.Error(http.StatusInternalServerError, "No branches found")
	}
	reqUser.Branches = branches

	return reqUser, nil
}

func (r *UserRepository) GetUserProfileByID(ctx context.Context, userID int32) (*externalUserDomain.User, error) {
	var currentUser externalUserDomain.User

	if err := r.SharedDbConnection.
		Unscoped().
		Where("id = ?", userID).
		First(&currentUser).Error; err != nil {
		return nil, err
	}

	var role *externalRoleDomain.Role
	if err := r.GetTenantDBConnection(ctx).Model(&role).Where("id = ?", currentUser.RoleID).First(&role).Error; err != nil {
		return nil, err
	}
	currentUser.Role = role

	return &currentUser, nil
}
