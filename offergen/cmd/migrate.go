package cmd

import (
	"offergen/config"
	"offergen/migrations"
)

type MigrateCmd struct{}

func (sc *MigrateCmd) Execute() {
	config := config.NewMigrateConfig()
	logger.Info("starting migrations")
	migrations.NewMigrator().Migrate(config.PostgresURL)
	logger.Info("finished migrations")
}
