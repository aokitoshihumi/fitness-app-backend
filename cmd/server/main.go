package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aokitoshihumi/fitness-app-backend/internal/config"
	"github.com/aokitoshihumi/fitness-app-backend/internal/database"
	"github.com/aokitoshihumi/fitness-app-backend/internal/handler"
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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.Healthz)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Ctrl+C / SIGTERM が来たら ctx が終わる
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

//サーバーが止まるまで待ち、止まった理由（エラー）を
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh: // エラーチャンネルにエラーが届く
		return err // ポート使用中など、起動に失敗した
	case <-ctx.Done(): // 終了シグナルがきた時
		log.Println("終了の合図を受け取りました")
	}

	// 処理中のリクエストを最大10秒待ってから止める
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("サーバーを停止しました")
	return nil
}
