// 一回分のワークアウトを入れる処理を書く
// workoutsテーブル ←その日の筋トレ（親）。誰が・いつ始めて・いつ終わったか
// workout_exercisesテーブル ←その日にやった1種目（子）。例: ベンチプレス、懸垂
// workout_setsテーブル ←その種目の1セット（孫）。例: 60kg × 10回
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aokitoshihumi/fitness-app-backend/internal/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ハンドラがエラーの種類でステータスコードを変えられるように、名前を付けておく
var (
	// この ID の記録が、ほかのユーザーのものとしてすでにある（409）
	ErrWorkoutIDConflict = errors.New("この ID の記録はほかのユーザーのものです")
	// 存在しない種目・部位を指定した、などで DB の制約に引っかかった（400）
	ErrInvalidWorkout = errors.New("記録の内容が正しくありません")
)

type WorkoutRepository struct {
	pool *pgxpool.Pool
}

func NewWorkoutRepository(pool *pgxpool.Pool) *WorkoutRepository {
	return &WorkoutRepository{pool: pool}
}

// Put は記録を保存する。同じ ID の記録があれば消してから入れ直す。
// 新しく作ったら created=true、置き換えたら false。
func (r *WorkoutRepository) Put(ctx context.Context, firebaseUID string, rec api.WorkoutRecord) (created bool, err error) {
	// トランザクション開始。ここから Commit までの変更は「全部成功」か「全部取り消し」になる
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("トランザクションを開始できませんでした: %w", err)
	}
	// 関数抜ける前にロールバック
	defer tx.Rollback(ctx)

	// Firebase の UID から users.id を探す
	var userID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE firebase_uid = $1`, firebaseUID).Scan(&userID)
	if err != nil {
		// pgx.ErrNoRows なら PUT /users/me がまだ呼ばれていない
		return false, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// 同じ ID の記録がすでにあるか、あれば誰のものかを調べる
	// rec.Id: 一回のワークアウト
	var ownerID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT user_id FROM workouts WHERE id = $1 FOR UPDATE`, rec.Id).Scan(&ownerID)
	switch {
	case errors.Is(err, pgx.ErrNoRows): // まだない場合
		created = true // 新規
	case err != nil: // その他エラーの時
		return false, fmt.Errorf("記録の確認に失敗しました: %w", err)
	case ownerID != userID: // リクエストしてきたIDと、取り出したユーザIDがちがかったら
		return false, ErrWorkoutIDConflict
	case ownerID == userID:
		// inser, updateなど結果を返さないクエリをトランザクションないで実行
		if _, err := tx.Exec(ctx, `DELETE FROM workouts WHERE id = $1`, rec.Id); err != nil {
			return false, fmt.Errorf("古い記録の削除に失敗しました: %w", err)
		}
	}

	// 親 → 子 → 孫 の順に入れる（外部キーがあるので親が先）
	_, err = tx.Exec(ctx, `
		INSERT INTO workouts (id, user_id, started_at, ended_at)
		VALUES ($1, $2, $3, $4)`,
		rec.Id, userID, rec.StartedAt, rec.EndedAt,
	)
	if err != nil {
		return false, wrapInsertError("記録", err)
	}

	// i が並び順（position）になる
	for i, ex := range rec.Exercises {
		_, err = tx.Exec(ctx, `
			INSERT INTO workout_exercises (id, workout_id, position, exercise_id, exercise_name, body_part)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			ex.Id, rec.Id, i, ex.Exercise.Id, ex.Exercise.Name, string(ex.Exercise.BodyPart),
		)
		if err != nil {
			return false, wrapInsertError("種目", err)
		}

		for j, set := range ex.Sets {
			// memo は省略できる（nil）。そのときは空文字で保存する
			memo := ""
			if set.Memo != nil {
				memo = *set.Memo
			}
			// weight は *float64 のまま渡す。nil なら NULL として保存される
			_, err = tx.Exec(ctx, `
				INSERT INTO workout_sets (id, workout_exercise_id, position, weight, reps, memo)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				set.Id, ex.Id, j, set.Weight, set.Reps, memo,
			)
			if err != nil {
				return false, wrapInsertError("セット", err)
			}
		}
	}

	// 4. ここまで全部成功したら確定する
	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("記録の保存を確定できませんでした: %w", err)
	}
	return created, nil
}

// wrapInsertError は、利用者の入力が原因の DB エラーを ErrInvalidWorkout にまとめる
func wrapInsertError(what string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503", // 外部キー違反：存在しない種目 ID や部位
			"23505", // 重複：同じ記録の中で同じ UUID を2回使った、など
			"23514": // CHECK 違反：終了日時が開始日時より前、など
			return fmt.Errorf("%sの保存に失敗しました: %w (%s)", what, ErrInvalidWorkout, pgErr.Message)
		}
	}
	return fmt.Errorf("%sの保存に失敗しました: %w", what, err)
}
