package chain

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"brahma/common/config"
	"brahma/common/logger"
	"brahma/internal/chain/erc20"
	"brahma/internal/chain/uniswap"
	"brahma/internal/models"
)

type UniswapPoolManager struct {
	ctx            context.Context
	log            logger.Logger
	client         *ethclient.Client
	blockRate      int
	contractStores map[common.Address]*contractStore
}

func NewUniswapPoolManager(ctx context.Context, log logger.Logger, config *config.ChainConfig) (*UniswapPoolManager, error) {
	client, err := ethclient.Dial(config.RpcUrl)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	contractStores := map[common.Address]*contractStore{}
	for _, poolAddress := range config.Pools {
		poolAddress := common.HexToAddress(poolAddress)
		poolContract, err := uniswap.NewUniswap(poolAddress, client)
		if err != nil {
			return nil, err
		}

		token0Address, err := poolContract.Token0(opts)
		if err != nil {
			return nil, err
		}

		token0Contract, err := erc20.NewErc20(token0Address, client)
		if err != nil {
			return nil, err
		}

		token1Address, err := poolContract.Token1(opts)
		if err != nil {
			return nil, err
		}

		token1Contract, err := erc20.NewErc20(token1Address, client)
		if err != nil {
			return nil, err
		}

		contractStores[poolAddress] = &contractStore{
			poolContract:   poolContract,
			token0Contract: token0Contract,
			token1Contract: token1Contract,
		}

		log.Info("Added new pool contract", "pool", poolAddress.Hex(), "token0", token0Address.Hex(), "token1", token1Address.Hex())
	}

	return &UniswapPoolManager{
		ctx:            ctx,
		log:            log,
		client:         client,
		blockRate:      config.BlockRate,
		contractStores: contractStores,
	}, nil
}

func (c *UniswapPoolManager) SubscribeBlocks(latestKnownBlock int64) (chan *models.PoolSnapshot, error) {
	headerCh := make(chan *types.Header)
	sub, err := c.client.SubscribeNewHead(c.ctx, headerCh)
	if err != nil {
		return nil, err
	}

	snapshotCh := make(chan *models.PoolSnapshot)
	go func() {
		c.log.Info("Starting tracking blocks", "block_rate", c.blockRate, "pool_count", len(c.contractStores))
		for {
			select {
			case err := <-sub.Err():
				// TODO: handle error
				c.log.Error("Block subscription channel killed", "err", err)
				break
			case header := <-headerCh:
				blockNumber := header.Number.Int64()
				if blockNumber < latestKnownBlock+int64(c.blockRate) {
					c.log.Debug("Skipping", "block", blockNumber)
					continue
				}

				c.log.Info("Taking snapshots", "block_number", blockNumber)

				for poolAddress := range c.contractStores {
					snapshot, err := c.getPoolSnapshot(poolAddress)
					if err != nil {
						// TODO: handle error
						c.log.Error("Unable to get pool snapshot", "pool_address", poolAddress.Hex(), "err", err)
						continue
					}

					snapshot.BlockNumber = blockNumber
					snapshotCh <- snapshot
				}

				latestKnownBlock = blockNumber
				c.log.Info("Finished taking snapshots")
			}
		}
	}()

	return snapshotCh, nil
}

func (c *UniswapPoolManager) getPoolSnapshot(poolAddress common.Address) (*models.PoolSnapshot, error) {
	opts := &bind.CallOpts{
		Pending: false,
		Context: c.ctx,
	}

	contractStore := c.contractStores[poolAddress]

	token0Balance, err := contractStore.token0Contract.BalanceOf(opts, poolAddress)
	if err != nil {
		return nil, err
	}

	token1Balance, err := contractStore.token1Contract.BalanceOf(opts, poolAddress)
	if err != nil {
		return nil, err
	}

	slot0, err := contractStore.poolContract.Slot0(opts)
	if err != nil {
		return nil, err
	}

	return &models.PoolSnapshot{
		PoolId:        poolAddress.Hex(),
		Token0Balance: token0Balance.Int64(),
		Token1Balance: token1Balance.Int64(),
		Tick:          slot0.Tick.Int64(),
	}, nil
}

type contractStore struct {
	poolContract   *uniswap.Uniswap
	token0Contract *erc20.Erc20
	token1Contract *erc20.Erc20
}
