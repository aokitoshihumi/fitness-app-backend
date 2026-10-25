package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port string // HTTPサーバが待ち受けるポート

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	FirebaseProjectID string // FirebaseのプロジェクトのID
}

// 環境変数を読み込んで Configを返す。
func Load() (*Config, error) {
	var missing []string

	// 一回の起動で足りないやつ全部返す
	cfg := &Config{
		Port: getEnv("PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     requireEnv("DB_USER", &missing),
		DBPassword: requireEnv("DB_PASSWORD", &missing),
		DBName:     requireEnv("DB_NAME", &missing),

		FirebaseProjectID: requireEnv("FIREBASE_PROJECT_ID", &missing),
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("必須の環境変数が設定されていません: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}

// 環境変数のkeyの値を渡す。
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// 環境変数のkeyの値を渡す。からならmissingに追加
func requireEnv(key string, missing *[]string) string {
	v := os.Getenv(key)
	if v == "" {
		*missing = append(*missing, key)
	}
	return v
}
