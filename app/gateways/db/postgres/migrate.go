package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func Migrate(databaseUrl string) error {
	driver, err := iofs.New(fs, "migrations")

	if err != nil {
		return fmt.Errorf("could not get embeded migration files: %s", err)
	}

	migration, err := migrate.NewWithSourceInstance("iofs", driver, databaseUrl)

	if err != nil {
		return fmt.Errorf("could not read migration files: %s", err)
	}

	err = migration.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {

		return err
	}

	return nil
}
