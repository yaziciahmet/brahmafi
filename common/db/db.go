package db

import (
	"context"

	"brahmafi/common/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	ctx    context.Context
	pool   *pgxpool.Pool
	config *config.DatabaseConfig
}

func NewDatabase(ctx context.Context, config *config.DatabaseConfig) *Database {
	return &Database{
		ctx:    ctx,
		config: config,
	}
}

func (d *Database) Connect() error {
	pool, err := pgxpool.New(d.ctx, d.config.DbUrl)
	if err != nil {
		return err
	}

	if err = pool.Ping(d.ctx); err != nil {
		return err
	}

	d.pool = pool

	return nil
}

func (d *Database) Migrate(migrationDir string) error {
	m, err := migrate.New(d.config.MigrationsDir, d.config.DbUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		switch err {
		case migrate.ErrNoChange:
			return nil
		case migrate.ErrNilVersion:
			return nil
		default:
			return err
		}
	}

	return nil
}

func (d *Database) GetPool() *pgxpool.Pool {
	return d.pool
}
