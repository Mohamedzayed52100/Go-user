package database

import (
	"fmt"

	"github.com/goplaceapp/goplace-user/internal/role/infrastructure/database/seeder"
)

func SyncSeeder(tx seeder.Executable) error {
	fmt.Println("Syncing seeders...")
	var err error

	err = seeder.UserDepartmentsSeeder(tx)
	if err != nil {
		fmt.Println("Error syncing seed: User departments")
		return err
	}
	fmt.Println("Synced Seed: User departments")

	err = seeder.PermissionsSeeder(tx)
	if err != nil {
		fmt.Println("Error syncing seed: Permissions")
		return err
	}
	fmt.Println("Synced Seed: Permissions")

	err = seeder.RolesSeeder(tx)
	if err != nil {
		fmt.Println("Error syncing seed: Roles")
		return err
	}
	fmt.Println("Synced Seed: Roles")

	err = seeder.RolePermissionAssignmentsSeeder(tx)
	if err != nil {
		fmt.Println("Error syncing seed: Role permission assignments")
		return err
	}
	fmt.Println("Synced Seed: Role permission assignments")

	return nil
}
