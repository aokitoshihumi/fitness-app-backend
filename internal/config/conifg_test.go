package config

import (
	"strings"
	"testing"
)

func TestLoad_必須がそろっていれば読み込める(t *testing.T) {
	t.Setenv("DB_USER", "kintore")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "kintore_db")
	t.Setenv("FIREBASE_PROJECT_ID", "kintore-dev") 
	t.Setenv("PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("エラーが出ないはずが出た: %v", err)
	}

	if cfg.DBUser != "kintore" {
		t.Errorf("DBUser = %q, want %q", cfg.DBUser, "kintore")
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q（デフォルト値）", cfg.Port, "8080")
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost = %q, want %q（デフォルト値）", cfg.DBHost, "localhost")
	}
	if cfg.FirebaseProjectID != "kintore-dev" {
		t.Errorf("FirebaseProjectID = %q, want %q", cfg.FirebaseProjectID, "kintore-dev")
	}
}

func TestLoad_必須が足りないとエラー(t *testing.T) {
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "kintore_db")

	_, err := Load()
	if err == nil {
		t.Fatal("エラーになるはずがならなかった")
	}

	// 足りない2つが両方メッセージに入っているか
	for _, key := range []string{"DB_USER", "DB_PASSWORD"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("エラーに %s が含まれていない: %v", key, err)
		}
	}
}
