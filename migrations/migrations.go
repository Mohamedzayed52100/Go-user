package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	shared "github.com/goplaceapp/goplace-user/migrations/shared"
	tenant "github.com/goplaceapp/goplace-user/migrations/tenant"
)

var SharedMigrations = []dbhelper.SqlxMigration{
	shared.CreateTenantsTable(),
	shared.CreateTenantProfilesTable(),
	shared.CreateDemoRequestsTable(),
	shared.CreateUsersTable(),
	shared.AlterUsersTableIndex(),
	shared.AlterUsersTableScheme(),
	shared.CreateTenantCredentialsTable(),
	shared.AddSetupCompletedToUsersTable(),
	shared.AddGenderToUsersTable(),
	shared.AlterTenantProfilesTable(),
	shared.CreateUserOtpsTable(),
	shared.AlterUserOtpsAddVerified(),
	shared.AlterUserOtpsAddTypeColumn(),
	shared.AlterUserOtpsAddVerifyMethodColumn(),
}

var TenantMigrations = []dbhelper.SqlxMigration{
	tenant.CreateBranchesTable(),
	tenant.CreateUserSessionsTable(),
	tenant.CreatePermissionsTable(),
	tenant.CreateUserDepartmentsTable(),
	tenant.CreateRolesTable(),
	tenant.CreateRolePermissionAssignmentsTable(),
	tenant.CreateUserBranchAssignmentsTable(),
	tenant.CreateBranchCredentialsTable(),
	tenant.AlterBranchesAddAdditionalData(),
	tenant.AlterRolesRemoveBranchId(),
	tenant.AlterPermissionsRemoveBranchId(),
	tenant.AlterUserDepartmentsRemoveBranchId(),
	tenant.AlterRolesTableAddConstraint(),
	tenant.AlterUserBranchAssignments(),
}
