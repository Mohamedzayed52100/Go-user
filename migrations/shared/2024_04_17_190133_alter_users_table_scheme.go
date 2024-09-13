package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterUsersTableScheme() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190133_alter_users_table_scheme",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			ALTER TABLE users
			DROP CONSTRAINT IF EXISTS users_email_ukey,
			DROP CONSTRAINT IF EXISTS users_phone_number_ukey;
		
			CREATE UNIQUE INDEX IF NOT EXISTS users_phone_number_tenant_id_uindex ON users (phone_number, tenant_id) WHERE phone_number != '' AND deleted_at IS NULL;
			CREATE UNIQUE INDEX IF NOT EXISTS users_email_tenant_id_uindex ON users (email, tenant_id) WHERE email != '' AND deleted_at IS NULL;
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
