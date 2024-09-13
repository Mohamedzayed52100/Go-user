package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterPermissionsRemoveBranchId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_12_145933_alter_permissions_remove_branch_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE permissions
            DROP COLUMN IF EXISTS branch_id;
			
			ALTER TABLE permissions
            DROP CONSTRAINT IF EXISTS permissions_name_unique;

			ALTER TABLE permissions
			ADD CONSTRAINT permissions_name_unique UNIQUE (name);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
