package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// サーバーを実際に起動しなくても、httptest を使えばハンドラだけを試せる。
//   - httptest.NewRequest … 偽のリクエストを作る
//   - httptest.NewRecorder … w の代わりに返事を記録してくれる入れ物
func TestHealthz_200とstatus_okを返す(t *testing.T) {
	// 手順:
	// 1. req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	// 2. rec := httptest.NewRecorder()
	rec := httptest.NewRecorder()
	// 3. Healthz(rec, req)            ← ハンドラを直接呼ぶ
	Healthz(rec, req)
	// 4. rec.Code が http.StatusOK か確かめる
	//      違ったら t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	// 5. rec.Header().Get("Content-Type") が "application/json" か確かめる
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	// 6. rec.Body.String() が `{"status":"ok"}` か確かめる
	if got := rec.Body.String(); got != `{"status":"ok"}` {
		t.Errorf("body = %q, want %q", got, `{"status":"ok"}`)
	}
	//
	// import に "net/http" と "net/http/httptest" が必要になる。
	// VSCode なら保存時に自動で追加される（されなければ手で足す）。
	//
	// 確認: go test ./internal/handler/ -v
	//   B1-1〜3 を書く前に実行すると FAIL、書いた後は PASS になるはず。
	//   「先にテストが落ちるのを見る」と、テストが本当に効いているとわかる。
}
