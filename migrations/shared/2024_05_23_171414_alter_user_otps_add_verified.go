package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func AlterUserOtpsAddVerified() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_23_171415_alter_user_otps_add_verified",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
            ALTER TABLE user_otps
            ADD COLUMN IF NOT EXISTS verified BOOLEAN DEFAULT FALSE,
			DROP CONSTRAINT IF EXISTS user_otps_token_ukey,
			ADD CONSTRAINT user_otps_token_ukey UNIQUE (user_id, token, verified);
            `
			_, err := tx.Exec(query)
			return err
		},
	}
}
