package domain

import "time"

type Branch struct {
	ID            int       `db:"id"`
	Name          string    `db:"name"`
	Country       string    `db:"country"`
	City          string    `db:"city"`
	Address       string    `db:"address"`
	GmapsLink     string    `db:"gmaps_link"`
	Email         string    `db:"email"`
	PhoneNumber   string    `db:"phone_number"`
	Website       string    `db:"website"`
	VatPercent    float32   `db:"vat_percent"`
	ServiceCharge float32   `db:"service_charge"`
	CrNumber      string    `db:"cr_number"`
	VatRegNumber  string    `db:"vat_reg_number"`
	Currency      string    `gorm:"-"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type BranchCredential struct {
	ID           int    `db:"id"`
	BranchID     int    `db:"branch_id"`
	ClientID     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
	Name         string `db:"name"`
	Enabled      bool   `db:"enabled"`
}

type UserBranchAssignment struct {
	ID       int32 `db:"id"`
	UserID   int32 `db:"user_id"`
	BranchID int32 `db:"branch_id"`
}
