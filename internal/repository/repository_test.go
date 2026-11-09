package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// newTestPool はテスト用の DB につなぎ、中身を空にしてから返す。
// TEST_DATABASE_URL が無ければ、そのテストはスキップする（DB なしでも go test ./... が通るように）
func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper() // 失敗したとき、この関数ではなく呼び出した行を表示してもらう

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL が無いので DB を使うテストはスキップします（make test-db で流せます）")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("DB の接続設定が不正です: %v", err)
	}
	// テストが終わったら接続を閉じる（defer のテスト版）
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("テスト用の DB につながりません: %v", err)
	}

	// 前のテストの残りを消す。users を消すと、ON DELETE CASCADE の先（workouts → exercises → sets）も消える。
	// 部位・種目のマスタ（シード）は消さない
	if _, err := pool.Exec(ctx, `TRUNCATE users CASCADE`); err != nil {
		t.Fatalf("テーブルを空にできませんでした: %v", err)
	}
	return pool
}

// countRows は表の行数を数える（カスケード削除の確認用）
func countRows(t *testing.T, pool *pgxpool.Pool, table string) int {
	t.Helper()
	var n int
	// 表の名前は $1 にできないので文字列でつなぐ。テストの中で決め打ちの名前しか渡さないので大丈夫
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM `+table).Scan(&n); err != nil {
		t.Fatalf("%s の行数を数えられませんでした: %v", table, err)
	}
	return n
}
