package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/aokitoshihumi/fitness-app-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 疎通確認の最大時間
const pingTimeout = 5 * time.Second

func Connect(ctx context.Context, cfg *config.Config)(*pgxpool.Pool, error){
	// poolはプールのポインタ
	pool, err := pgxpool.New(ctx, buildDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("DBの接続設定が不正です: %w", err)
	}

	// 一回疎通する
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	// 必ずcancelを呼ぶ
	defer cancel()

	if err := pool.Ping(ctx);
	err != nil {
		pool.Close()
		return nil, fmt.Errorf("DB(%s:%s/%s)に接続できません: %w", cfg.DBHost, cfg.DBPort, cfg.DBName, err)
	}

	return pool, nil
}

// 設定から接続文字を組み立てる
func buildDSN(cfg *config.Config) string {
	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}
	return u.String()
}