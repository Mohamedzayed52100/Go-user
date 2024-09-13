package domain

import "time"
type TenantCredential struct {
	ID           int    `db:"id"`
	TenantID     string `db:"tenant_id"`
	ClientID     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
	Name         string `db:"name"`
	Enabled      bool   `db:"enabled"`
}

type TenantProfile struct {
	ID          string `db:"id"`
	TenantID    string `db:"tenant_id"`
	Name        string `db:"name"`
	DisplayName string `db:"display_name"`
	Address     string `db:"address"`
	Phone       string `db:"phone"`
	Email       string `db:"email"`
	Logo        string `db:"logo"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
