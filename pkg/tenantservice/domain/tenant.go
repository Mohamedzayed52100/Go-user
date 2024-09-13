package domain

import "time"

type Tenant struct {
	ID        string    `db:"id"`
	Domain    string    `db:"domain"`
	DbName    string    `db:"db_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
