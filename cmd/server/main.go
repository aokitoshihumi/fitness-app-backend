package main

import (
	"log"
	"net/http"
	"os"
)

func main(){
	// Cloud Run は待ち受けるポートを環境変数 PORT で渡してくる。
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz) // ヘルスチェック

	log.Printf("listening on :%s", port)
	//サーバーを起動する。もし失敗して戻ってきたら理由を出力して終了
	log.Fatal(http.ListenAndServe(":"+port, mux)) 
}

// サーバー生きているかの確認エンドポイント
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // OKのステータスコード200を返す
	w.Write([]byte("healthy")) // healthyという文字を返す
}