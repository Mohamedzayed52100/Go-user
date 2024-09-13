package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterRolesRemoveBranchId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_12_144929_alter_roles_remove_branch_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE roles 
            DROP COLUMN IF EXISTS branch_id,
			DROP CONSTRAINT IF EXISTS roles_branch_id_fk;
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
