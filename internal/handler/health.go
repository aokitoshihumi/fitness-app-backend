package handler

import (
	"net/http"
)

// Healthz は「サーバーが生きているか」を返すハンドラ。
// Cloud Run や監視ツールが定期的に叩いて、200 が返ってくれば正常とみなす。
//
// 返したいレスポンス:
//
//	ステータス: 200 OK
//	ヘッダー:   Content-Type: application/json
//	本文:       {"status":"ok"}
//
//   - w … 返事を書き込む先（ここにヘッダー・ステータス・本文を書く）
//   - r … 届いたリクエスト（今回は使わない）
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`)) // Writeはバイトしか受け取らない
}
