package domain

import "time"

type Permission struct {
	ID          int32     `db:"id"`
	Roles       []int32   `gorm:"-"`
	Name        string    `db:"name"`
	DisplayName string    `db:"display_name"`
	Description string    `db:"description"`
	Category    string    `db:"category"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
