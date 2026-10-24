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
	// Ctrl+C / SIGTERM が来たら ctx が終わる。
	// 一番最初に作って、この ctx を DB 接続などにも渡す。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		return err
	}
	// run を抜ける時に接続を閉じる
	defer pool.Close()

	log.Printf("DBに接続しました (%s:%s/%s)", cfg.DBHost, cfg.DBPort, cfg.DBName)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Healthz)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// サーバーを裏で動かし、止まった理由（エラー）を errCh で受け取る
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh: // ポート使用中など、起動に失敗した
		return err
	case <-ctx.Done(): // 終了の合図が来た
		log.Println("終了の合図を受け取りました")
	}

	// 処理中のリクエストを最大8秒待ってから止める
	// （Cloud Run は SIGTERM の約10秒後に強制終了するので、それより短くする）
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("サーバーを停止しました")
	return nil
}
