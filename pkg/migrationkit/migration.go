package migrationkit

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"github.com/amikai/line-go-prac/pkg/logkit"
)

type MigrationConfig struct {
	// see https://github.com/golang-migrate/migrate
	Source string
	URL    string
}

type Migration struct {
	*migrate.Migrate
	logger *logkit.Logger
}

func (m *Migration) Up() error {
	err := m.Migrate.Up()
	if err == nil {
		return nil
	}

	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info("no change to migrate")
		return nil
	}

	return err
}

func (m *Migration) Close() error {
	serr, derr := m.Migrate.Close()
	if serr != nil {
		return serr
	}

	if derr != nil {
		return derr
	}

	return nil
}

func NewMigration(ctx context.Context, conf *MigrationConfig) *Migration {
	logger := logkit.FromContext(ctx).With(
		zap.String("source", conf.Source),
		zap.String("url", conf.URL),
	)

	m, err := migrate.New(conf.Source, conf.URL)
	if err != nil {
		logger.Fatal("failed to create migration", zap.Error(err))
	}

	logger.Info("create migration successfully")

	return &Migration{
		Migrate: m,
		logger:  logger,
	}
}
