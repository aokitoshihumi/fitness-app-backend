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
	tests := []struct {
		name   string
		header string
	}{
		{"Bearer", "Bearer good-token"},
		{"小文字のbearer", "bearer good-token"},
		{"大文字のBEARER", "BEARER good-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Middleware(fakeVerifier{})(http.HandlerFunc(echoUID))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.header)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			if got := rec.Body.String(); got != "user-123" {
				t.Errorf("body = %q, want %q", got, "user-123")
			}
		})
	}
}

func TestMiddleware_だめなときは401(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"ヘッダーがない", ""},
		{"Bearer がない", "good-token"},
		{"Bearer だけでスペースがない", "Bearer"},
		{"トークンが空", "Bearer "},
		{"Bearer 以外の方式", "Basic good-token"},
		{"トークンが不正", "Bearer bad-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 次のハンドラが呼ばれたら記録する。401 のときは呼ばれてはいけない。
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})
			h := Middleware(fakeVerifier{})(next)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if called {
				t.Error("認証に失敗したのに次のハンドラが呼ばれた")
			}
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
			if got := rec.Header().Get("WWW-Authenticate"); got != "Bearer" {
				t.Errorf("WWW-Authenticate = %q, want %q", got, "Bearer")
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json")
			}
			want := `{"error":{"code":"unauthorized","message":"認証に失敗しました"}}`
			if got := rec.Body.String(); got != want {
				t.Errorf("body = %q, want %q", got, want)
			}
		})
	}
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{"Bearer", "Bearer abc", "abc", true},
		{"小文字", "bearer abc", "abc", true},
		{"JWTの形でもそのまま取り出せる", "Bearer aaa.bbb.ccc", "aaa.bbb.ccc", true},
		{"空文字", "", "", false},
		{"スペースがない", "Bearer", "", false},
		{"トークンが空", "Bearer ", "", false},
		{"別の方式", "Basic abc", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotOK := bearerToken(tt.header)
			if gotToken != tt.wantToken || gotOK != tt.wantOK {
				t.Errorf("bearerToken(%q) = (%q, %v), want (%q, %v)",
					tt.header, gotToken, gotOK, tt.wantToken, tt.wantOK)
			}
		})
	}
}
