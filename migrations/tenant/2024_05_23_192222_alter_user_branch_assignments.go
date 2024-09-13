package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterUserBranchAssignments() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_23_192222_alter_user_branch_assignments",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE user_branch_assignments
            ADD CONSTRAINT uk_user_branch_assignments UNIQUE (user_id, branch_id);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
