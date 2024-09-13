package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateRolesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190719_create_roles_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS roles (
					id SERIAL PRIMARY KEY,
					department_id INTEGER NOT NULL,
					name TEXT NOT NULL,
					display_name TEXT NOT NULL,
					branch_id INTEGER,
					created_at 	timestamptz DEFAULT NOW(),
					updated_at 	timestamptz DEFAULT NOW(),

					CONSTRAINT roles_branch_id_fk FOREIGN KEY (branch_id) REFERENCES branches (id) ON DELETE CASCADE,
					CONSTRAINT fk_department_id FOREIGN KEY(department_id) REFERENCES user_departments(id)
				)
				`
			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE roles`
			_, err := tx.Exec(query)
			return err
		},
	}
}
