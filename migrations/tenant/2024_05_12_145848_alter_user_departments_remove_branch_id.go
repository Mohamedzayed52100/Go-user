package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterUserDepartmentsRemoveBranchId() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_12_145849_alter_user_departments_remove_branch_id",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE user_departments
            DROP COLUMN IF EXISTS branch_id;

			ALTER TABLE user_departments
            DROP CONSTRAINT IF EXISTS user_departments_name_branch_id_ukey;

			ALTER TABLE user_departments
            ADD CONSTRAINT user_departments_name_branch_id_ukey UNIQUE (name);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
