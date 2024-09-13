package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateBranchCredentialsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_191156_create_branch_credentials_table",
		Migrate: func(tx *sqlx.Tx) error {
			// Create the branch_credentials table
			createTableQuery := `
			CREATE TABLE IF NOT EXISTS branch_credentials (
				id SERIAL PRIMARY KEY,
				branch_id INTEGER NOT NULL,
				name text NOT NULL,
				client_id text NOT NULL,
				client_secret text NOT NULL,
				enabled boolean NOT NULL default true,
				created_at timestamp with time zone NOT NULL DEFAULT now(),
				updated_at timestamp with time zone NOT NULL DEFAULT now(),

				CONSTRAINT branch_credentials_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
			)`
			_, err := tx.Exec(createTableQuery)
			if err != nil {
				return err
			}

			// Check if client_id and client_secret columns exist in the branches table
			checkColumnsQuery := `
			SELECT column_name
			FROM information_schema.columns
			WHERE table_name='branches' AND (column_name='client_id' OR column_name='client_secret')
			`

			rows, err := tx.Query(checkColumnsQuery)
			if err != nil {
				return err
			}

			var columns []string
			for rows.Next() {
				var columnName string
				if err := rows.Scan(&columnName); err != nil {
					return err
				}
				columns = append(columns, columnName)
			}

			if len(columns) < 2 {
				return nil
			}

			// Create a temporary table to store client_id and client_secret
			createTempTableQuery := `
			CREATE TEMPORARY TABLE IF NOT EXISTS temp_branches AS 
		   	SELECT client_id, client_secret FROM branches
   			`
			_, err = tx.Exec(createTempTableQuery)
			if err != nil {
				return err
			}

			// Insert the client_id and client_secret from the temporary table into the branch_credentials table
			insertQuery := `
   			INSERT INTO branch_credentials (client_id, client_secret)
   			SELECT client_id, client_secret FROM temp_branches
   			`
			_, err = tx.Exec(insertQuery)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS branch_credentials`
			_, err := tx.Exec(query)
			return err
		},
	}
}
