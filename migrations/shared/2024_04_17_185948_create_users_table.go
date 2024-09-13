package migrations

import (
	"os"
	"time"

	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"
	externalUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"github.com/jmoiron/sqlx"
	"github.com/tinygg/gofaker"
	"golang.org/x/crypto/bcrypt"
)

func CreateUsersTable() dbhelper.SqlxMigration {
	return dbhelper.SqlxMigration{
		ID: "2024_04_17_185948_create_users_table",
		Migrate: func(tx *sqlx.Tx) error {
			query := `
CREATE TABLE IF NOT EXISTS users (
	id 			SERIAL NOT NULL,
	employee_id text UNIQUE NULL,
	first_name	text,
	last_name		text,
	email 		text,
	password 		text NOT NULL,
	role_id integer NOT NULL,
	phone_number 	text NULL,
	birthdate     timestamp NOT NULL,
	avatar 		text NULL,
	status 		text DEFAULT 'active',
	pin_code 		integer NULL,
	tenant_id 	uuid NOT NULL,
	branch_id integer NULL,
	joined_at timestamp NOT NULL,
	created_at 	timestamptz DEFAULT NOW(),
	updated_at 	timestamptz DEFAULT NOW(),
	deleted_at timestamp NULL,

	PRIMARY KEY ("id"),
	CONSTRAINT "users_email_ukey" UNIQUE ("email"),
	CONSTRAINT "users_phone_number_ukey" UNIQUE ("phone_number"),
	CONSTRAINT "users_tenant_id_fkey" FOREIGN KEY ("tenant_id") REFERENCES tenants(id) ON DELETE CASCADE
);
`
			_, err := tx.Exec(query)
			return err
		},
		Seed: func(tx *sqlx.Tx) error {
			environment := os.Getenv("ENVIRONMENT")
			if environment != meta.ProdEnvironment && environment != meta.StagingEnvironment {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456789"), bcrypt.DefaultCost)
				if err != nil {
					return err
				}

				var tenants []domain.Tenant
				if err := tx.Select(&tenants, `SELECT id FROM tenants`); err != nil {
					return err
				}

				for _, tenant := range tenants {
					users := []externalUserDomain.User{
						{
							FirstName:   "John",
							EmployeeID:  "1",
							LastName:    "Doe",
							Email:       "john@goplace.io",
							Password:    string(hashedPassword),
							Birthdate:   gofaker.Date(),
							PhoneNumber: gofaker.Phone(),
							BranchID:    1,
							RoleID:      1,
							PinCode:     "12345678",
							TenantID:    tenant.ID,
							JoinedAt:    time.Now(),
						},
					}

					var count int
					if err := tx.Get(&count, `SELECT COUNT(*) FROM users WHERE tenant_id = $1`, tenant.ID); err != nil {
						return err
					}
					if count == 0 {
						if _, err := tx.NamedExec(`INSERT INTO users (employee_id, birthdate, joined_at, first_name, last_name, email, password, phone_number, branch_id, role_id, pin_code, tenant_id) VALUES (:employee_id, :birthdate, :joined_at, :first_name, :last_name, :email, :password, :phone_number, :branch_id, :role_id, :pin_code, :tenant_id)`, users); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
		Rollback: func(tx *sqlx.Tx) error {
			query := `DROP TABLE users;`
			_, err := tx.Exec(query)
			return err
		},
	}
}
