package database

import (
	"testing"

	"github.com/aokitoshihumi/fitness-app-backend/internal/config"
)

func TestBuild(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "kintore",
		DBPassword: "p@ss:word/1",
		DBName:     "kintore_db",
		DBSSLMode:  "require",
	}

	got := buildDSN(cfg)
	want := "postgres://kintore:p%40ss%3Aword%2F1@localhost:5432/kintore_db?sslmode=require"
	if got != want {
		t.Errorf("buildDSN() = %q, want %q", got, want)
	}
}

func TestNewPoolConfig_接続数が絞られている(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "kintore",
		DBPassword: "secret",
		DBName:     "kintore_db",
		DBSSLMode:  "disable",
	}

	poolCfg, err := newPoolConfig(cfg)
	if err != nil {
		t.Fatalf("エラーが出ないはずが出た: %v", err)
	}
	if poolCfg.MaxConns != 3 {
		t.Errorf("MaxConns = %d, want %d", poolCfg.MaxConns, 3)
	}
}

func TestNewPoolConfig_sslmodeが不正だとエラー(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "kintore",
		DBPassword: "secret",
		DBName:     "kintore_db",
		DBSSLMode:  "yes",
	}

	if _, err := newPoolConfig(cfg); err == nil {
		t.Fatal("エラーになるはずがならなかった")
	}
}
