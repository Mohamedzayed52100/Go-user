package domain

import "time"

type UserFeedback struct {
	ID          int32     `db:"id"`
	UserID      int32     `db:"user_id"`
	Rate        float32   `db:"rate"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
