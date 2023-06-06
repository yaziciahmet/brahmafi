package main

import (
	"context"
	"fmt"
	"os"

	"brahmafi/common/config"
	"brahmafi/common/db"
	"brahmafi/common/logger"
)

func main() {
	config, err := config.NewConfiguration()
	if err != nil {
		fmt.Printf("failed to load configuration, %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	log := logger.NewLogger(config.GetLoggerConfig())

	db := db.NewDatabase(ctx, config.GetDatabaseConfig(), log.Clone("db"))
	if err = db.Connect(); err != nil {
		log.Fatal("failed to connect to database", "err", err)
	}

	if err = db.Migrate(); err != nil {
		log.Fatal("failed to run migrations", "err", err)
	}
}
