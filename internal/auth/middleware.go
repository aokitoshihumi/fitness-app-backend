package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type TokenVerifaier interface {
	VerifyIDToken(ctx context.Context, idToken string)(string, error)
}

type ctxKey struct{}

// 最初にトークンを検証する
// 正しければ uid を context に入れる
// 次の処理に渡す
func Middleware(v TokenVerifaier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idToken, ok:= strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || idToken == "" {
				writeUnauthorized(w)
				return
			}

			uid, err := v.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				log.Printf("ID トークンの検証に失敗しました: %v", err)
				writeUnauthorized(w)
				return 
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UIDFromContext(ctx context.Context)(string,bool) {
	uid, ok := ctx.Value(ctxKey{}).(string)
	return uid, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":{"code":"unauthorized","message":"認証に失敗しました"}}`))
}