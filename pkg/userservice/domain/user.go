package domain

import (
	"time"

	externalRoleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"gorm.io/gorm"
)

type User struct {
	ID             int32                    `db:"id"`
	EmployeeID     string                   `db:"employee_id"`
	FirstName      string                   `db:"first_name"`
	LastName       string                   `db:"last_name"`
	Email          string                   `db:"email"`
	Password       string                   `db:"password"`
	PhoneNumber    string                   `db:"phone_number"`
	Birthdate      time.Time                `db:"birthdate"`
	Gender         string                   `db:"gender"`
	RoleID         int32                    `db:"role_id"`
	Role           *externalRoleDomain.Role `gorm:"-"`
	Avatar         string                   `db:"avatar"`
	Status         string                   `db:"status"`
	PinCode        string                   `db:"pin_code"`
	TenantID       string                   `db:"tenant_id"`
	BranchID       int32                    `db:"branch_id"`
	Branch         *Branch                  `gorm:"-"`
	Branches       []*Branch                `gorm:"-"`
	Timezone       string                   `gorm:"-"`
	SetupCompleted bool                     `db:"setup_completed"`
	JoinedAt       time.Time                `db:"joined_at"`
	CreatedAt      time.Time                `db:"created_at"`
	UpdatedAt      time.Time                `db:"updated_at"`
	DeletedAt      gorm.DeletedAt
}

type UserSession struct {
	UserID     int32      `db:"user_id"`
	SessionID  string     `db:"session_id"`
	LoginTime  time.Time  `db:"login_time"`
	LogoutTime *time.Time `db:"logout_time"`
}
