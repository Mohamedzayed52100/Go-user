package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AddGenderToUsersTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_04_150822_add_gender_to_users_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE users ADD COLUMN gender VARCHAR(255) DEFAULT 'Male' NOT NULL;
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "ALTER TABLE users DROP COLUMN gender"
			_, err := tx.Exec(query)
			return err
		},
	}
}
