package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateUserSessionsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190628_create_user_sessions_table",
		Migrate: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`
                CREATE TABLE IF NOT EXISTS user_sessions (
                    user_id INTEGER NOT NULL,
					session_id TEXT NOT NULL,
					login_time TIMESTAMP NOT NULL,
					logout_time TIMESTAMP NULL,

					CONSTRAINT pk_user_sessions PRIMARY KEY (user_id, session_id)
				);
            `)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			_, err := tx.Exec(`DROP TABLE IF EXISTS user_sessions`)
			return err
		},
	}
}
