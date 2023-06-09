package core

import (
	"brahmafi/common/logger"
	"brahmafi/internal/chain"
	"brahmafi/internal/models"
	"brahmafi/internal/repository"
	"strconv"
)

type BrahmaService struct {
	log         logger.Logger
	repository  *repository.PoolRepository
	poolManager *chain.UniswapPoolManager
}

func NewBrahmaService(log logger.Logger, repository *repository.PoolRepository, poolManager *chain.UniswapPoolManager) *BrahmaService {
	return &BrahmaService{
		log:         log,
		repository:  repository,
		poolManager: poolManager,
	}
}

func (b *BrahmaService) WatchBlocks() {
	ch, err := b.poolManager.SubscribeBlocks(0)
	if err != nil {
		b.log.Fatal("Failed to subscribe to blocks", "err", err)
	}

	for snapshot := range ch {
		b.log.Info("Got", "snapshot", snapshot)
		if err := b.repository.CreatePoolSnapshot(snapshot); err != nil {
			b.log.Error("Failed to insert snapshot into repository", "err", err)
		}
	}
}

func (b *BrahmaService) GetPoolSnapshotByBlock(poolId, blockNumber string) (*models.PoolSnapshot, error) {
	if blockNumber == "latest" {
		snapshot, err := b.repository.GetPoolSnapshotOfLatestBlock(poolId)
		if err != nil {
			return nil, err
		}

		return snapshot, nil
	}

	block, err := strconv.Atoi(blockNumber)
	if err != nil {
		return nil, err
	}

	snapshot, err := b.repository.GetPoolSnapshotNearestToBlock(poolId, block)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (b *BrahmaService) GetPoolSnapshots(poolId string) ([]*models.PoolSnapshotWithTokenDelta, error) {
	return b.repository.GetPoolSnapshotHistory(poolId)
}

func (b *BrahmaService) AddPool(poolAddress string) error {
	return b.repository.CreatePool(&models.Pool{Address: poolAddress})
}
