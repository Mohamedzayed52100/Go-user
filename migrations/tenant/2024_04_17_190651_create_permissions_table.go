package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreatePermissionsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190652_create_permissions_table",
		Migrate: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`
                CREATE TABLE IF NOT EXISTS permissions(
                    id SERIAL PRIMARY KEY,
                    name TEXT NOT NULL,
                    display_name TEXT NOT NULL,
					description TEXT NULL,
					category TEXT NULL,
					branch_id INTEGER NOT NULL,
                    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

					FOREIGN KEY (branch_id) REFERENCES branches (id) ON DELETE CASCADE,
					CONSTRAINT permissions_name_unique UNIQUE (name, branch_id)
                )
            `)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS permissions`)
			return err
		},
	}
}
