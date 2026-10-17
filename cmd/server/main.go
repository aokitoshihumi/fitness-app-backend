package main

import (
	"context"
	"log"
	"os"

	"github.com/aokitoshihumi/fitness-app-backend/internal/config"
	"github.com/aokitoshihumi/fitness-app-backend/internal/database"
)

	func main() {
	if err := run(); err != nil {
		log.Printf("起動に失敗しました: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// contextを作る
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		return err
	}

	// runを抜ける時に接続を閉じる
	defer pool.Close()

	log.Printf("DBに接続しました (%s:%s/%s)", cfg.DBHost, cfg.DBPort, cfg.DBName)
	return nil
}