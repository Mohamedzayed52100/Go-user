package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/jmoiron/sqlx"
)

func CreateUserBranchAssignmentsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_191930_create_user_branch_assignments_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS user_branch_assignments (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL,
				branch_id INTEGER NOT NULL,

				CONSTRAINT fk_branches FOREIGN KEY (branch_id) REFERENCES branches(id)
			)
			`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				assignment := &externalUserDomain.UserBranchAssignment{
					UserID:   1,
					BranchID: 1,
				}

				_, err := tx.NamedExec("INSERT INTO user_branch_assignments(user_id, branch_id) VALUES (:user_id, :branch_id)", assignment)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS user_branch_assignments`
			_, err := tx.Exec(query)
			return err
		}}
}
