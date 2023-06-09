package main

import (
	"context"
	"fmt"
	"os"

	"brahma/common/config"
	"brahma/common/db"
	"brahma/common/logger"
	"brahma/internal/api"
	"brahma/internal/chain"
	"brahma/internal/core"
	"brahma/internal/repository"
)

func main() {
	config, err := config.NewConfiguration()
	if err != nil {
		fmt.Printf("Failed to load configuration, %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	log := logger.NewLogger(config.GetLoggerConfig())

	db := db.NewDatabase(ctx, config.GetDatabaseConfig(), log.Clone("db"))
	if err = db.Connect(); err != nil {
		log.Fatal("Failed to connect to database", "err", err)
	}

	if err = db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations", "err", err)
	}

	poolRepository := repository.NewPoolRepository(ctx, db)

	uniswapPoolManager, err := chain.NewUniswapPoolManager(ctx, log.Clone("pool_manager"), config.GetChainConfig())
	if err != nil {
		log.Fatal("Failed to create uniswap pool manager", "err", err)
	}

	brahmaService := core.NewBrahmaService(log.Clone("core"), poolRepository, uniswapPoolManager)
	// save passed pool addresses into database if not already exists
	for _, poolAddress := range config.GetChainConfig().Pools {
		brahmaService.AddPool(poolAddress)
	}

	go brahmaService.WatchBlocks()

	server := api.NewServer(log, brahmaService, config.GetApiConfig())
	server.Listen()
}
