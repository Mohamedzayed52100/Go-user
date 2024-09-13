package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterUserOtpsAddVerifyMethodColumn() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_23_171424_alter_user_otps_add_verify_method_column",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE user_otps
            ADD COLUMN IF NOT EXISTS verify_method VARCHAR(255) DEFAULT 'email';

			ALTER TABLE user_otps
			DROP CONSTRAINT IF EXISTS user_otps_token_ukey,
			ADD CONSTRAINT user_otps_token_ukey UNIQUE (user_id, token, verified, type, verify_method);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
