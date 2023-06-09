package repository

import (
	"context"

	"brahma/common/db"
	"brahma/internal/models"

	"github.com/jackc/pgx/v5"
)

const (
	createPoolQuery = `
		INSERT INTO pool (address)
		VALUES ($1)
		ON CONFLICT DO NOTHING
	`
	createPoolSnapshotQuery = `
		INSERT INTO pool_snapshot (pool_id, token0_balance, token1_balance, tick, block_number)
		VALUES ($1, $2, $3, $4, $5)
	`
	getPoolSnapshotOfLatestBlockQuery = `
		SELECT id, pool_id, token0_balance, token1_balance, tick, block_number, taken_at
		FROM pool_snapshot
		WHERE pool_id = $1
		ORDER BY block_number DESC
		LIMIT 1
	`
	getPoolSnapshotNearestToBlockQuery = `
		SELECT id, pool_id, token0_balance, token1_balance, tick, block_number, taken_at
		FROM pool_snapshot
		WHERE abs($1 - block_number) = (
			SELECT min(abs($1 - block_number))
			FROM pool_snapshot
			WHERE pool_id = $2
		) AND pool_id = $2
		ORDER BY block_number DESC
		LIMIT 1
	`
	getPoolSnapshotHistoryQuery = `
		SELECT 
			id,
			pool_id,
			token0_balance,
			token1_balance,
			tick,
			block_number,
			taken_at,
			token0_balance - (lead(token0_balance) OVER (ORDER BY block_number DESC)) AS "token0_delta",
			token1_balance - (lead(token1_balance) OVER (ORDER BY block_number DESC)) AS "token1_delta"
		FROM pool_snapshot
		WHERE pool_id = $1
	`
)

type PoolRepository struct {
	ctx context.Context
	db  *db.Database
}

func NewPoolRepository(ctx context.Context, db *db.Database) *PoolRepository {
	return &PoolRepository{
		ctx: ctx,
		db:  db,
	}
}

func (p *PoolRepository) CreatePool(pool *models.Pool) error {
	_, err := p.db.GetPool().Exec(p.ctx, createPoolQuery, pool.Address)
	return err
}

func (p *PoolRepository) CreatePoolSnapshot(snapshot *models.PoolSnapshot) error {
	_, err := p.db.GetPool().Exec(
		p.ctx,
		createPoolSnapshotQuery,
		snapshot.PoolId,
		snapshot.Token0Balance,
		snapshot.Token1Balance,
		snapshot.Tick,
		snapshot.BlockNumber,
	)
	return err
}

func (p *PoolRepository) GetPoolSnapshotOfLatestBlock(poolId string) (*models.PoolSnapshot, error) {
	row := p.db.GetPool().QueryRow(p.ctx, getPoolSnapshotOfLatestBlockQuery, poolId)
	return p.scanPoolSnapshot(row)
}

func (p *PoolRepository) GetPoolSnapshotNearestToBlock(poolId string, blockNumber int) (*models.PoolSnapshot, error) {
	row := p.db.GetPool().QueryRow(p.ctx, getPoolSnapshotNearestToBlockQuery, blockNumber, poolId)
	return p.scanPoolSnapshot(row)
}

func (p *PoolRepository) GetPoolSnapshotHistory(poolId string) ([]*models.PoolSnapshotWithTokenDelta, error) {
	rows, err := p.db.GetPool().Query(p.ctx, getPoolSnapshotHistoryQuery, poolId)
	if err != nil {
		return nil, err
	}

	var snapshots []*models.PoolSnapshotWithTokenDelta
	for rows.Next() {
		snapshot := &models.PoolSnapshotWithTokenDelta{}

		err = rows.Scan(
			&snapshot.Id,
			&snapshot.PoolId,
			&snapshot.Token0Balance,
			&snapshot.Token1Balance,
			&snapshot.Tick,
			&snapshot.BlockNumber,
			&snapshot.TakenAt,
			&snapshot.Token0Delta,
			&snapshot.Token1Delta,
		)
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, snapshot)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return snapshots, nil
}

func (p *PoolRepository) scanPoolSnapshot(row pgx.Row) (*models.PoolSnapshot, error) {
	snapshot := &models.PoolSnapshot{}
	err := row.Scan(
		&snapshot.Id,
		&snapshot.PoolId,
		&snapshot.Token0Balance,
		&snapshot.Token1Balance,
		&snapshot.Tick,
		&snapshot.BlockNumber,
		&snapshot.TakenAt,
	)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
