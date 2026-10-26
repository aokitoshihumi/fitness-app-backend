package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func (r *UserRepository) CreateIfNotExists(ctx context.Context, firebaseUID string) (created bool, err error) {
	// 行を返さない SQL なので QueryRow ではなく Exec を使う
	// tagには何行のデータが変更されたのかという情報が入る。後どんな操作だったのかも
	tag, err := r.pool.Exec(ctx, `
		INSERT INTO users (firebase_uid) VALUES ($1)
		ON CONFLICT (firebase_uid) DO NOTHING`,
		firebaseUID,
	)
	if err != nil {
		return false, fmt.Errorf("ユーザーの登録に失敗しました: %w", err)
	}
	// 入れた行数。作ったら 1、すでにいて何もしなかったら 0
	return tag.RowsAffected() == 1, nil
}
