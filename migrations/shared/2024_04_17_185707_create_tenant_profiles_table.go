package migrations

import (
	"os"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/internal/tenant/domain"
	tenantDomain "github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
)

func CreateTenantProfilesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_185707_create_tenant_profiles_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
CREATE TABLE IF NOT EXISTS tenant_profiles (
id uuid NOT NULL DEFAULT uuid_generate_v4(),
tenant_id uuid NOT NULL,
"name" text NOT NULL,
address text NULL,
phone text NOT NULL,
email text NOT NULL,
logo text NULL,
created_at timestamptz DEFAULT now(),
updated_at timestamptz DEFAULT now(),

PRIMARY KEY ("id"),
FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
)`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				var tenants []tenantDomain.Tenant
				if err := tx.Select(&tenants, `SELECT * FROM tenants`); err != nil {
					return err
				}

				for _, tenant := range tenants {
					tenantProfile := domain.TenantProfile{
						TenantID: tenant.ID,
						Name:     "starbucks_eg",
						Address:  gofaker.City(),
						Phone:    gofaker.Phone(),
						Email:    gofaker.Email(),
						Logo:     "https://placehold.co/250x150",
					}
					if _, err := tx.NamedExec(`INSERT INTO tenant_profiles (tenant_id, name, address, phone, email, logo) VALUES (:tenant_id, :name, :address, :phone, :email, :logo)`, tenantProfile); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS tenant_profiles;`
			_, err := tx.Exec(query)
			return err
		},
	}
}
