package domain

import "time"

type UserOtp struct {
	ID           int32     `db:"id"`
	UserID       int32     `db:"user_id"`
	User         *User     `db:"user"`
	Code         string    `db:"code"`
	Token        string    `db:"token"`
	Type         string    `db:"type"`
	Verified     bool      `db:"verified"`
	VerifyMethod string    `db:"verify_method"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
