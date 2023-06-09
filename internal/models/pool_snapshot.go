package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolSnapshot struct {
	Id            uuid.UUID `json:"id"`
	PoolId        string    `json:"poolId"`
	Token0Balance int64     `json:"token0Balance"`
	Token1Balance int64     `json:"token1Balance"`
	Tick          int64     `json:"tick"`
	BlockNumber   int64     `json:"blockNumber"`
	TakenAt       time.Time `json:"takenAt"`
}

type PoolSnapshotWithTokenDelta struct {
	PoolSnapshot
	Token0Delta *int64 `json:"token0Delta"`
	Token1Delta *int64 `json:"token1Delta"`
}
