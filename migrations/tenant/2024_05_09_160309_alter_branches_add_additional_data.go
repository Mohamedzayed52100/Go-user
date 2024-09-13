package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterBranchesAddAdditionalData() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_09_160118_alter_branches_add_additional_data",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE branches
            ADD COLUMN IF NOT EXISTS city TEXT NULL,
            ADD COLUMN IF NOT EXISTS address TEXT NULL,
            ADD COLUMN IF NOT EXISTS gmaps_link TEXT NULL,
            ADD COLUMN IF NOT EXISTS email TEXT NULL,
            ADD COLUMN IF NOT EXISTS phone_number TEXT NULL,
            ADD COLUMN IF NOT EXISTS website TEXT NULL
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
