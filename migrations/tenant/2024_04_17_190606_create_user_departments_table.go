package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateUserDepartmentsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190607_create_user_departments_table",
		Migrate: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`
                CREATE TABLE IF NOT EXISTS user_departments (
                    id SERIAL PRIMARY KEY,
                    name TEXT NOT NULL,
					branch_id INTEGER,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW(),

					CONSTRAINT user_departments_name_branch_id_ukey UNIQUE (name, branch_id),
					CONSTRAINT user_departments_branch_id_fk FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
				)
            `)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS user_departments`)
			return err
		},
	}
}
