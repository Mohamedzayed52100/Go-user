package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateUserOtpsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_05_19_160358_create_user_otps_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			CREATE TABLE IF NOT EXISTS user_otps (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL,
				token TEXT NOT NULL,
				code TEXT NOT NULL,
				expires_at TIMESTAMP NOT NULL,
                created_at TIMESTAMPTZ DEFAULT NOW(),
                updated_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT user_otps_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				CONSTRAINT user_otps_token_ukey UNIQUE (user_id)
			)
			`

			_, err := tx.Exec(query)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `
			DROP TABLE IF EXISTS user_otps
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}
