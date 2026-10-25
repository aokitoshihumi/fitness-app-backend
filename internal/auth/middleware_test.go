package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeVerifier は偽物として用意
// "good-token" だけを本物として扱い、UID "user-123" を返す。
type fakeVerifier struct{}

func (fakeVerifier) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	if idToken == "good-token" {
		return "user-123", nil
	}
	return "", errors.New("invalid token")
}

// ミドルウェアの先にあるハンドラの代わり。context の UID を本文に書き出す。
func echoUID(w http.ResponseWriter, r *http.Request) {
	uid, ok := UIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(uid))
}

func TestMiddleware_正しいトークンならUIDが次に渡る(t *testing.T) {
	h := Middleware(fakeVerifier{})(http.HandlerFunc(echoUID))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "user-123" {
		t.Errorf("body = %q, want %q", got, "user-123")
	}
}

func TestMiddleware_だめなときは401(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"ヘッダーがない", ""},
		{"Bearer がない", "good-token"},
		{"トークンが空", "Bearer "},
		{"トークンが不正", "Bearer bad-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Middleware(fakeVerifier{})(http.HandlerFunc(echoUID))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
			want := `{"error":{"code":"unauthorized","message":"認証に失敗しました"}}`
			if got := rec.Body.String(); got != want {
				t.Errorf("body = %q, want %q", got, want)
			}
		})
	}
}
