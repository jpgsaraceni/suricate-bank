package postgres

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Migrate(databaseUrl, sourceUrl string) error {
	migration, err := migrate.New(
		sourceUrl,
		databaseUrl)

	if err != nil {
		return fmt.Errorf("could not read migration files: %s", err)
	}

	err = migration.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {

		return err
	}

	return nil
}
