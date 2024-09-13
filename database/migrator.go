package database

import (
	"github.com/goplaceapp/goplace-common/pkg/dbhelper"
	"github.com/goplaceapp/goplace-user/migrations"
	"gorm.io/gorm"
)

type Migrator interface {
	SharedMigrationUp(db *gorm.DB) error
	TenantMigrationUp(db *gorm.DB) error
}

func SharedMigrationUp(db *gorm.DB) error {
	sqlMigration := dbhelper.Sqlx{
		Migrations: migrations.SharedMigrations,
	}

	conn, _ := db.DB()
	err := sqlMigration.Migrate(conn, "postgres")
	if err != nil {
		return err
	}

	return nil
}

func TenantMigrationUp(db *gorm.DB) error {
	sqlMigration := dbhelper.Sqlx{
		Migrations: migrations.TenantMigrations,
	}

	conn, _ := db.DB()
	err := sqlMigration.Migrate(conn, "postgres")
	if err != nil {
		return err
	}

	if err := SyncSeeder(conn); err != nil {
		return err
	}

	return nil
}
