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
	}

	got := buildDSN(cfg)
	want := "postgres://kintore:p%40ss%3Aword%2F1@localhost:5432/kintore_db"
	if got != want {
		t.Errorf("buildDSN() = %q, want %q", got, want)
	}
}
