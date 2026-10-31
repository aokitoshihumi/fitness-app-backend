package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/aokitoshihumi/fitness-app-backend/internal/auth"
)

type UserCreator interface {
	CreateIfNotExists(ctx context.Context, firebaseUID string) (created bool, err error)
}

type UsersHandler struct {
	users UserCreator
}

func NewUsersHandler(users UserCreator) *UsersHandler {
	return &UsersHandler{users: users}
}

// PUT /api/v1/users/meのハンドラ
// auth.Middlewareの後ろに置く
// 新規作成は201
// すでにある場合は204
func (h *UsersHandler) PutMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UIDFromContext(r.Context())
	if !ok {
		// ミドルウェアを通っていれば必ず入ってくる。
		log.Println("PUT /users/me: contextにUIDがありません。auth.Middlewareを通していますか")
		writeInternalError(w)
		return
	}

	created, err := h.users.CreateIfNotExists(r.Context(), uid)
	if err != nil {
		log.Printf("PUT /users/me %v", err)
		writeInternalError(w)
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error":{"code":"internal_error","message":"サーバーでエラーが発生しました"}}`))
}
