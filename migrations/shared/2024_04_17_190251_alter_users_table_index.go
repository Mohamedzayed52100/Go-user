package migrations

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/jmoiron/sqlx"
)

func AlterUsersTableIndex() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_190251_alter_users_table_index",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
			DROP INDEX IF EXISTS users_phone_number_tenant_id_uindex;
			DROP INDEX IF EXISTS users_email_tenant_id_uindex;
			`
			_, err := tx.Exec(query)
			if err != nil {
				return err
			}

			// Get all users employeeIds
			result, err := tx.Query(`SELECT id, employee_id FROM users;`)
			if err != nil {
				return err
			}
			defer result.Close()

			var users []externalUserDomain.User
			for result.Next() {
				var id int32
				var employeeId string
				err := result.Scan(&id, &employeeId)
				if err != nil {
					return err
				}
				users = append(users, externalUserDomain.User{ID: id, EmployeeID: employeeId})
			}

			_, err = tx.Exec(`ALTER TABLE users DROP COLUMN employee_id;`)
			if err != nil {
				return err
			}

			_, err = tx.Exec(`ALTER TABLE users ADD COLUMN employee_id TEXT NULL;`)
			if err != nil {
				return err
			}

			for _, user := range users {
				tx.NamedExec(`UPDATE users SET employee_id = :employee_id WHERE id = :id`, user)
			}

			return nil
		},
	}
}
