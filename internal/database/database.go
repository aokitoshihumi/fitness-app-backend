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

// 1プロセスが同時に持つDBの最大接続数
// Lambdaの実行数×maxConnsがRDSのmax_connectionsを超えないようにする
const maxConns = 3

func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := newPoolConfig(cfg)
	if err != nil {
		return nil, err
	}

	// poolはプールのポインタ
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("DBの接続プールを作れません: %w", err)
	}

	// 一回疎通する
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	// 必ずcancelを呼ぶ
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("DB(%s:%s/%s)に接続できません: %w", cfg.DBHost, cfg.DBPort, cfg.DBName, err)
	}

	return pool, nil
}

// 接続文字列を読み込み、プールの設定（接続数など）を足して返す
func newPoolConfig(cfg *config.Config) (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig(buildDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("DBの接続設定が不正です: %w", err)
	}
	poolCfg.MaxConns = maxConns
	return poolCfg, nil
}

// 設定から接続文字を組み立てる
func buildDSN(cfg *config.Config) string {
	// sslmode を書かないと pgx は prefer（暗号化できなければ暗号化なし）で動くので、必ず付ける
	q := url.Values{}
	q.Set("sslmode", cfg.DBSSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:     net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:     cfg.DBName,
		RawQuery: q.Encode(),
	}
	return u.String()
}
