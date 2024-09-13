package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AddSetupCompletedToUsersTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_02_124531_add_setup_completed_to_users_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE users
            
            ADD COLUMN setup_completed BOOLEAN DEFAULT FALSE
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := "ALTER TABLE users DROP COLUMN setup_completed"
			_, err := tx.Exec(query)
			return err
		},
	}
}
