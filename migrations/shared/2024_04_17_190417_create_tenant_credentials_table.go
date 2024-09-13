package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/jmoiron/sqlx"
)

func CreateTenantCredentialsTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190417_create_tenant_credentials_table",
		Migrate: func(tx *sqlx.Tx) error {
			// Create the tenant_credentials table
			createTableQuery := `
			CREATE TABLE IF NOT EXISTS tenant_credentials (
                id SERIAL PRIMARY KEY,
                tenant_id uuid NOT NULL,
    			name text NOT NULL,
                client_id text NOT NULL,
                client_secret text NOT NULL,
    			enabled boolean NOT NULL default true,
				created_at timestamp with time zone NOT NULL DEFAULT now(),
				updated_at timestamp with time zone NOT NULL DEFAULT now(),

				CONSTRAINT tenant_credentials_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
			)`
			_, err := tx.Exec(createTableQuery)
			if err != nil {
				return err
			}

			// Check if client_id and client_secret columns exist in the tenants table
			checkColumnsQuery := `
			SELECT column_name 
			FROM information_schema.columns 
			WHERE table_name='tenants' AND (column_name='client_id' OR column_name='client_secret')
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
			CREATE TEMPORARY TABLE IF NOT EXISTS temp_tenants AS
			SELECT client_id, client_secret FROM tenants
			`
			_, err = tx.Exec(createTempTableQuery)
			if err != nil {
				return err
			}

			// Insert the client_id and client_secret from the temporary table into the tenant_credentials table
			insertQuery := `
   			INSERT INTO tenant_credentials (client_id, client_secret)
   			SELECT client_id, client_secret FROM temp_tenants
   			WHERE client_id IS NOT NULL AND client_secret IS NOT NULL
   			`
			_, err = tx.Exec(insertQuery)
			return err
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE IF EXISTS tenant_credentials`
			_, err := tx.Exec(query)
			return err
		},
	}
}
