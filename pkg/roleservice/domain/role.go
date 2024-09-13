package domain

import (
	"time"

	"github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"
)

type Role struct {
	ID           int32                  `db:"id"`
	DepartmentID int32                  `db:"department_id"`
	Department   *domain.UserDepartment `gorm:"-"`
	Name         string                 `db:"name"`
	DisplayName  string                 `db:"display_name"`
	Permissions  []*Permission          `gorm:"-"`
	UsersCount   int32                  `gorm:"-"`
	CreatedAt    time.Time              `db:"created_at"`
	UpdatedAt    time.Time              `db:"updated_at"`
}
