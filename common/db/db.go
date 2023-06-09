package db

import (
	"context"

	"brahma/common/config"
	"brahma/common/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	ctx    context.Context
	pool   *pgxpool.Pool
	config *config.DatabaseConfig
	log    logger.Logger
}

func NewDatabase(ctx context.Context, config *config.DatabaseConfig, log logger.Logger) *Database {
	return &Database{
		ctx:    ctx,
		config: config,
		log:    log,
	}
}

func (d *Database) Connect() error {
	d.log.Info("Connecting to database")

	pool, err := pgxpool.New(d.ctx, d.config.DbUrl)
	if err != nil {
		return err
	}

	if err = pool.Ping(d.ctx); err != nil {
		return err
	}

	d.pool = pool

	d.log.Info("Successfully connected to database", "db_url", d.config.DbUrl)
	return nil
}

func (d *Database) Migrate() error {
	m, err := migrate.New(d.config.MigrationsDir, d.config.DbUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		switch err {
		case migrate.ErrNoChange:
			break
		case migrate.ErrNilVersion:
			break
		default:
			return err
		}
	}

	d.log.Info("Migrations ran successfully")
	return nil
}

// get pool connection manager
func (d *Database) GetPool() *pgxpool.Pool {
	return d.pool
}
