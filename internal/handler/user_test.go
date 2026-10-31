package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aokitoshihumi/fitness-app-backend/internal/auth"
)

// fakeUsers は DB の代わり。決めておいた値を返し、受け取った UID を覚えておく。
type fakeUsers struct {
	created bool
	err     error

	gotUID string
}

func (f *fakeUsers) CreateIfNotExists(ctx context.Context, firebaseUID string) (bool, error) {
	f.gotUID = firebaseUID
	return f.created, f.err
}

// 本物の auth.Middleware に偽の検証役を渡して、context に UID を入れてもらう
type fakeVerifier struct{}

func (fakeVerifier) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	return "user-123", nil
}

// ミドルウェア → PutMe の順に通して、返事を返す
func callPutMe(users UserCreator) *httptest.ResponseRecorder {
	h := auth.Middleware(fakeVerifier{})(http.HandlerFunc(NewUsersHandler(users).PutMe))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPutMe_新規なら201_既存なら204(t *testing.T) {
	tests := []struct {
		name       string
		created    bool
		wantStatus int
	}{
		{"新しく作った", true, http.StatusCreated},
		{"すでにいた", false, http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &fakeUsers{created: tt.created}

			rec := callPutMe(users)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if users.gotUID != "user-123" {
				t.Errorf("CreateIfNotExists に渡った UID = %q, want %q", users.gotUID, "user-123")
			}
			if rec.Body.Len() != 0 {
				t.Errorf("body = %q, want 空", rec.Body.String())
			}
		})
	}
}

func TestPutMe_DBエラーなら500(t *testing.T) {
	users := &fakeUsers{err: errors.New("db is down")}

	rec := callPutMe(users)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	want := `{"error":{"code":"internal_error","message":"サーバーでエラーが発生しました"}}`
	if got := rec.Body.String(); got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestPutMe_ミドルウェアを通っていなければ500(t *testing.T) {
	users := &fakeUsers{}
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()

	// ミドルウェアを通さずに直接呼ぶ
	NewUsersHandler(users).PutMe(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
