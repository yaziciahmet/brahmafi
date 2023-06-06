package models

import "github.com/google/uuid"

type Pool struct {
	Id        uuid.UUID
	Address   string
	ChainId   string
	ChainName string
}
