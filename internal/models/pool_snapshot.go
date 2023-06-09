package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolSnapshot struct {
	Id            uuid.UUID
	PoolId        string
	Token0Balance int64
	Token1Balance int64
	Tick          int64
	BlockNumber   int64
	TakenAt       time.Time
}

type PoolSnapshotWithTokenDelta struct {
	PoolSnapshot
	Token0Delta *int64
	Token1Delta *int64
}
