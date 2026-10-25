package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

type ctxKey struct{}

// 最初にトークンを検証する
// 正しければ uid を context に入れる
// 次の処理に渡す
func Middleware(v TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idToken, ok := bearerToken(r.Header.Get("Authorization"))
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

// Bearerに設定されているトークンを取り出す。bearerとは区別しない
func bearerToken(header string)(string, bool){
	scheme, token, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return token, true
}

func UIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(ctxKey{}).(string)
	return uid, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer") // RFC9110より、401の時はBearerヘッダーをつけることがMUSTになっているらしい
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":{"code":"unauthorized","message":"認証に失敗しました"}}`))
}
