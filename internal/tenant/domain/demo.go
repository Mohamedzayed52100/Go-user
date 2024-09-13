package domain

import "time"

type Demo struct {
	Id             string    `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	PhoneNumber    string    `db:"phone_number"`
	Country        string    `db:"country"`
	RestaurantName string    `db:"restaurant_name"`
	BranchesNo     int32     `db:"branches_no"`
	FirstTimeCrm   bool      `db:"first_time_crm"`
	SystemName     string    `db:"system_name"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
