package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterTenantProfilesTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_09_154226_alter_tenant_profiles_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE tenant_profiles
            ADD COLUMN IF NOT EXISTS display_name text NOT NULL DEFAULT '';
            `
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			return nil
		},
	}
}
