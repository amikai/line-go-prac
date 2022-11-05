package bot

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/amikai/line-go-prac/config"
	"github.com/amikai/line-go-prac/pkg/logkit"
	"github.com/amikai/line-go-prac/pkg/migrationkit"
)

func newMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migration",
		Short: "runs the messages migration job",
		RunE:  runMigration,
	}
}

func runMigration(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	conf, err := config.LoadConf()
	if err != nil {
		log.Fatal("failed to load config", err.Error())
	}

	var logger *logkit.Logger
	{
		loggerConfig := &logkit.LoggerConfig{
			Level:       logkit.LoggerLevel(conf.Logger.Level),
			Development: conf.Logger.Developement,
		}
		logger = logkit.NewLogger(loggerConfig)
	}
	defer func() {
		_ = logger.Sync()
	}()
	ctx = logger.WithContext(ctx)

	var migration *migrationkit.Migration
	{
		migrationConfig := &migrationkit.MigrationConfig{
			Source: conf.Migration.Source,
			URL:    conf.Migration.URL,
		}
		migration = migrationkit.NewMigration(ctx, migrationConfig)
		if err := migration.Up(); err != nil {
			logger.Fatal("failed to run migration", zap.Error(err))
		}

		logger.Info("run migration job successfully, terminating ...")
	}

	return nil
}
