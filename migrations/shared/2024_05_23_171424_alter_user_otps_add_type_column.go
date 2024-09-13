package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterUserOtpsAddTypeColumn() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_23_171424_alter_user_otps_add_type_column",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE user_otps
            ADD COLUMN IF NOT EXISTS type VARCHAR(255) DEFAULT 'reset_password' NOT NULL;

			ALTER TABLE user_otps
			DROP CONSTRAINT IF EXISTS user_otps_token_ukey,
			ADD CONSTRAINT user_otps_token_ukey UNIQUE (user_id, token, verified, type);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
