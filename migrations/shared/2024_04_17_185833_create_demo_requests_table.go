package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/internal/tenant/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
)

func CreateDemoRequestsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_185833_create_demo_requests_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
CREATE TABLE IF NOT EXISTS demo_requests (
	id UUID NOT NULL DEFAULT uuid_generate_v4(),
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	phone_number TEXT NOT NULL,
	country TEXT,
	restaurant_name TEXT NOT NULL,
	branches_no INTEGER NOT NULL,
	first_time_crm BOOLEAN,
	system_name TEXT,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),

	CONSTRAINT demo_requests_pkey PRIMARY KEY (id),
	CONSTRAINT demo_requests_email_ukey UNIQUE (email)
)
`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				demos := []domain.Demo{
					{
						Name:           gofaker.Name(),
						Email:          "bot@goplace.io",
						PhoneNumber:    gofaker.Phone(),
						Country:        gofaker.Country(),
						RestaurantName: gofaker.Company(),
						BranchesNo:     1,
						FirstTimeCrm:   true,
						SystemName:     "Goplace",
					},
				}

				for _, demo := range demos {
					var count int
					err := tx.Get(&count, `SELECT COUNT(*) FROM demo_requests WHERE email = $1`, demo.Email)
					if err != nil {
						return err
					}
					if count == 0 {
						if _, err := tx.NamedExec(`INSERT INTO demo_requests (name, email, phone_number, country, restaurant_name, branches_no, first_time_crm, system_name) VALUES (:name, :email, :phone_number, :country, :restaurant_name, :branches_no, :first_time_crm, :system_name)`, demo); err != nil {
							return err
						}
					}
				}
				return nil
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS demo_requests`
			_, err := tx.Exec(query)
			return err
		},
	}
}
