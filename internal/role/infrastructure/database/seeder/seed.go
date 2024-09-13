package seeder

import (
	"database/sql"
	"errors"
	departmentDomain "github.com/goplaceapp/goplace-user/pkg/departmentservice/domain"

	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/rbac"
)

type Executable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func UserDepartmentsSeeder(tx Executable) error {
	query := `
	INSERT INTO user_departments (name, created_at, updated_at)
	VALUES ($1, NOW(), NOW())
	ON CONFLICT (name)
	DO NOTHING
	`

	for _, department := range meta.Departments {
		_, err := tx.Exec(query, department.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDepartmentByName(tx Executable, name string) (*departmentDomain.UserDepartment, error) {
	query := `SELECT * FROM user_departments WHERE name = $1`
	row := tx.QueryRow(query, name)

	department := &departmentDomain.UserDepartment{}
	err := row.Scan(&department.ID, &department.Name, &department.CreatedAt, &department.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return department, nil
}

func RolesSeeder(tx Executable) error {
	query := `
    INSERT INTO roles (department_id, name, display_name, created_at, updated_at)
 	VALUES ($1, $2, $3, NOW(), NOW())
    ON CONFLICT (name, department_id)
    DO UPDATE SET
        department_id = $1,
        name = $2,
        display_name = $3,
        updated_at = NOW()
    `

	for _, r := range rbac.Roles {
		getManagementDepartment, err := GetDepartmentByName(tx, "Management")
		if err != nil {
			return err
		}

		_, err = tx.Exec(query, getManagementDepartment.ID, r.Name, r.DisplayName)
		if err != nil {
			return err
		}
	}

	return nil
}

func PermissionsSeeder(tx Executable) error {
	// Insert or update permissions
	insertQuery := `
        INSERT INTO permissions (name, display_name, description, category, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        ON CONFLICT (name)
        DO UPDATE SET display_name = $2, description = $3, category = $4, updated_at = NOW()
    `

	for _, permission := range rbac.AllPermissions {
		_, err := tx.Exec(insertQuery, permission.Name, permission.DisplayName, permission.Description, permission.Category)
		if err != nil {
			return err
		}
	}

	// Get all permissions from the database
	rows, err := tx.Query("SELECT name FROM permissions")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Delete permissions not in rbac.AllPermissions
	deleteQuery := "DELETE FROM permissions WHERE name = $1"
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}

		if !rbac.PermissionExists(name) {
			_, err := tx.Exec(deleteQuery, name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RolePermissionAssignmentsSeeder(tx Executable) error {
	roleRow := tx.QueryRow("SELECT id FROM roles WHERE name = 'super-admin'")
	var roleID int
	if err := roleRow.Scan(&roleID); err != nil {
		return err
	}

	// Get all permissions
	permissions, err := tx.Query("SELECT id FROM permissions")
	if err != nil {
		return err
	}
	defer permissions.Close()

	// Assign all permissions to the super admin role
	for permissions.Next() {
		var permissionID int
		if err := permissions.Scan(&permissionID); err != nil {
			return err
		}
		_, err := tx.Exec(`
                INSERT INTO role_permission_assignments (role_id, permission_id)
                VALUES ($1, $2)
                ON CONFLICT (role_id, permission_id)
                DO NOTHING`, roleID, permissionID)
		if err != nil {
			return err
		}
	}

	return nil
}
