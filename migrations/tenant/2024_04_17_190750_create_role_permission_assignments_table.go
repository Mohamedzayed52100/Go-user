package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateRolePermissionAssignmentsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190751_create_role_permission_assignments_table",
		Migrate: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`
                CREATE TABLE IF NOT EXISTS role_permission_assignments(
                    id SERIAL PRIMARY KEY,
                    role_id INT NOT NULL,
                    permission_id INT NOT NULL,

                    CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id),
                    CONSTRAINT fk_role_permission_assignments_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
                    CONSTRAINT fk_role_permission_assignments_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
                )
            `)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS role_permission_assignments`)
			return err
		},
	}
}
