package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/jmoiron/sqlx"
)

func CreateBranchesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190506_create_branches_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS branches (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				country TEXT NOT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			)
			`

			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				branches := []externalUserDomain.Branch{
					{
						Name:    "Branch 1",
						Country: "Saudi Arabia",
					},
				}
				for _, branch := range branches {
					_, err := tx.NamedExec(`INSERT INTO branches (name, country) VALUES (:name, :country)`, branch)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS branches`
			_, err := tx.Exec(query)
			return err
		},
	}
}
