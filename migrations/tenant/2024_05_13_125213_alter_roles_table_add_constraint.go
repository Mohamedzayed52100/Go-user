package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterRolesTableAddConstraint() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_13_125213_alter_roles_table_add_constraint",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE roles
            DROP CONSTRAINT IF EXISTS uk_department_name,
            ADD CONSTRAINT uk_department_name UNIQUE(department_id, name);
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := ""
			_, err := tx.Exec(query)
			return err
		},
	}
}
