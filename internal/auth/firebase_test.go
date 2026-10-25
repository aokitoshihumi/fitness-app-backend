package auth

import (
	"context"
	"testing"
)

func TestNewFirebaseVerifier_エミュレーターの環境変数があると起動しない(t *testing.T) {
	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")

	v, err := NewFirebaseVerifier(context.Background(), "kintore-dev")
	if err == nil {
		t.Fatal("エラーになるはずが、nil が返った")
	}
	if v != nil {
		t.Errorf("verifier = %v, want nil", v)
	}
}

func TestNewFirebaseVerifier_環境変数がなければ作れる(t *testing.T) {
	// 手元の環境に設定されていても、このテストの間だけ空にする
	t.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "")

	v, err := NewFirebaseVerifier(context.Background(), "kintore-dev")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if v == nil {
		t.Fatal("verifier が nil")
	}
}
