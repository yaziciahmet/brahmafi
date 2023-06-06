package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolSnapshot struct {
	Id            uuid.UUID
	PoolId        uuid.UUID
	Token0Balance int
	Token1Balance int
	Tick          int
	BlockNumber   int
	TakenAt       time.Time
}

type PoolSnapshotWithTokenDelta struct {
	PoolSnapshot
	Token0Delta *int
	Token1Delta *int
}
