package migrations

import (
	"os"

	"github.com/google/uuid"
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"
	"github.com/jmoiron/sqlx"
)

func CreateTenantsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_185536_create_tenants_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
CREATE TABLE IF NOT EXISTS tenants (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "domain" text NOT NULL,
    "db_name" text NOT NULL,
    "created_at" timestamptz DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),

    PRIMARY KEY ("id"),
    CONSTRAINT "tenants_domain_ukey" UNIQUE ("domain"),
    CONSTRAINT "tenants_db_name_ukey" UNIQUE ("db_name")
);`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				uuid := uuid.New()

				tenants := []domain.Tenant{
					{
						ID:     uuid.String(),
						Domain: "starbucks.com",
						DbName: "starbucks_eg",
					},
				}

				for _, tenant := range tenants {
					var count int
					err := tx.Get(&count, `SELECT COUNT(*) FROM tenants WHERE db_name = $1`, tenant.DbName)
					if err != nil {
						return err
					}
					if count == 0 {
						if _, err := tx.NamedExec(`INSERT INTO tenants (id, domain, db_name) VALUES (:id, :domain, :db_name)`, tenant); err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS "tenants"`
			_, err := tx.Exec(query)
			return err
		},
	}

}
