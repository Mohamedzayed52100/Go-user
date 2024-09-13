package domain

type RolePermissionAssignment struct {
	ID           int32 `db:"id"`
	RoleID       int32 `db:"role_id"`
	PermissionID int32 `db:"permission_id"`
}
