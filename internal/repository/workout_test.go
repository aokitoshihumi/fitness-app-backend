package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aokitoshihumi/fitness-app-backend/internal/api"
	"github.com/google/uuid"
)

const (
	alice = "firebase-uid-alice"
	bob   = "firebase-uid-bob"
)

// newRecord はテスト用の記録を作る。ベンチプレス2セット → スクワット1セット
func newRecord(startedAt time.Time) api.WorkoutRecord {
	weight := 60.0
	memo := "補助あり"
	return api.WorkoutRecord{
		Id:        uuid.New(),
		StartedAt: startedAt,
		EndedAt:   startedAt.Add(time.Hour),
		Exercises: []api.WorkoutExercise{
			{
				Id:       uuid.New(),
				Exercise: api.Exercise{Id: "bench-press", Name: "ベンチプレス", BodyPart: api.Chest},
				Sets: []api.WorkoutSet{
					{Id: uuid.New(), Weight: &weight, Reps: 10, Memo: &memo},
					{Id: uuid.New(), Weight: nil, Reps: 8}, // 重さなし・memo なし
				},
			},
			{
				Id:       uuid.New(),
				Exercise: api.Exercise{Id: "squat", Name: "スクワット", BodyPart: api.Legs},
				Sets: []api.WorkoutSet{
					{Id: uuid.New(), Weight: &weight, Reps: 5},
				},
			},
		},
	}
}

// setupRepos は DB を空にして、alice と bob を登録した状態のリポジトリを返す
func setupRepos(t *testing.T) *WorkoutRepository {
	t.Helper()
	pool := newTestPool(t)
	users := NewUserRepository(pool)
	for _, uid := range []string{alice, bob} {
		if _, err := users.CreateIfNotExists(context.Background(), uid); err != nil {
			t.Fatalf("ユーザー %s を作れませんでした: %v", uid, err)
		}
	}
	return NewWorkoutRepository(pool)
}

// 1. 新しく保存して、そのまま読み戻せる（並び順・重さなし・memo なしも含めて）
func TestWorkoutRepository_PutAndList(t *testing.T) {
	repo := setupRepos(t)
	ctx := context.Background()
	rec := newRecord(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC))

	created, err := repo.Put(ctx, alice, rec)
	if err != nil {
		t.Fatalf("Put でエラー: %v", err)
	}
	if !created {
		t.Errorf("新しい記録なのに created=false")
	}

	got, err := repo.List(ctx, alice, nil, nil)
	if err != nil {
		t.Fatalf("List でエラー: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("記録の数 = %d, 期待 1", len(got))
	}
	r := got[0]

	if r.Id != rec.Id {
		t.Errorf("id = %v, 期待 %v", r.Id, rec.Id)
	}
	// time は == ではなく Equal で比べる（タイムゾーンの持ち方が違っても同じ時刻なら true）
	if !r.StartedAt.Equal(rec.StartedAt) || !r.EndedAt.Equal(rec.EndedAt) {
		t.Errorf("日時 = %v〜%v, 期待 %v〜%v", r.StartedAt, r.EndedAt, rec.StartedAt, rec.EndedAt)
	}

	// 種目の並び順：ベンチプレス → スクワット
	if len(r.Exercises) != 2 {
		t.Fatalf("種目の数 = %d, 期待 2", len(r.Exercises))
	}
	if r.Exercises[0].Exercise.Id != "bench-press" || r.Exercises[1].Exercise.Id != "squat" {
		t.Errorf("種目の順番 = %s, %s, 期待 bench-press, squat",
			r.Exercises[0].Exercise.Id, r.Exercises[1].Exercise.Id)
	}

	// セットの並び順と中身
	sets := r.Exercises[0].Sets
	if len(sets) != 2 {
		t.Fatalf("ベンチプレスのセット数 = %d, 期待 2", len(sets))
	}
	if sets[0].Id != rec.Exercises[0].Sets[0].Id || sets[1].Id != rec.Exercises[0].Sets[1].Id {
		t.Errorf("セットの順番が保存したときと違います")
	}
	if sets[0].Weight == nil || *sets[0].Weight != 60 {
		t.Errorf("1セット目の weight = %v, 期待 60", sets[0].Weight)
	}
	if sets[1].Weight != nil {
		t.Errorf("2セット目の weight = %v, 期待 nil（NULL）", *sets[1].Weight)
	}
	if *sets[0].Memo != "補助あり" {
		t.Errorf("1セット目の memo = %q, 期待 %q", *sets[0].Memo, "補助あり")
	}
	if *sets[1].Memo != "" {
		t.Errorf("memo を送らなかったセットの memo = %q, 期待 空文字", *sets[1].Memo)
	}
}

// 2. 同じ ID でもう一度 PUT すると、丸ごと置き換わる
func TestWorkoutRepository_PutReplaces(t *testing.T) {
	repo := setupRepos(t)
	ctx := context.Background()
	rec := newRecord(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC))

	if _, err := repo.Put(ctx, alice, rec); err != nil {
		t.Fatalf("1回目の Put でエラー: %v", err)
	}

	// スクワットを消して、ベンチプレスのセットを1つにする
	rec.Exercises = rec.Exercises[:1]
	rec.Exercises[0].Sets = rec.Exercises[0].Sets[:1]

	created, err := repo.Put(ctx, alice, rec)
	if err != nil {
		t.Fatalf("2回目の Put でエラー: %v", err)
	}
	if created {
		t.Errorf("置き換えなのに created=true")
	}

	got, err := repo.List(ctx, alice, nil, nil)
	if err != nil {
		t.Fatalf("List でエラー: %v", err)
	}
	if len(got) != 1 || len(got[0].Exercises) != 1 || len(got[0].Exercises[0].Sets) != 1 {
		t.Fatalf("置き換え後の形が違います: %+v", got)
	}
	// 古い種目・セットが DB に残っていないこと
	if n := countRows(t, repo.pool, "workout_exercises"); n != 1 {
		t.Errorf("workout_exercises の行数 = %d, 期待 1", n)
	}
	if n := countRows(t, repo.pool, "workout_sets"); n != 1 {
		t.Errorf("workout_sets の行数 = %d, 期待 1", n)
	}
}

// 3. ほかの人の記録の ID で PUT すると ErrWorkoutIDConflict。元の記録は変わらない
func TestWorkoutRepository_PutOtherUsersID(t *testing.T) {
	repo := setupRepos(t)
	ctx := context.Background()
	rec := newRecord(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC))

	if _, err := repo.Put(ctx, alice, rec); err != nil {
		t.Fatalf("alice の Put でエラー: %v", err)
	}

	_, err := repo.Put(ctx, bob, rec)
	if !errors.Is(err, ErrWorkoutIDConflict) {
		t.Fatalf("bob の Put のエラー = %v, 期待 ErrWorkoutIDConflict", err)
	}

	// alice の記録はそのまま、bob には何もない
	aliceRecs, _ := repo.List(ctx, alice, nil, nil)
	bobRecs, _ := repo.List(ctx, bob, nil, nil)
	if len(aliceRecs) != 1 || len(bobRecs) != 0 {
		t.Errorf("記録の数 alice=%d bob=%d, 期待 1 と 0", len(aliceRecs), len(bobRecs))
	}
}

// 4. 存在しない種目・終了が開始より前、は ErrInvalidWorkout。途中まで保存されたりしない
func TestWorkoutRepository_PutInvalid(t *testing.T) {
	// テーブル駆動テスト：入力の作り方と名前を並べて、同じ確認を繰り返す
	tests := []struct {
		name   string
		modify func(rec *api.WorkoutRecord)
	}{
		{
			name: "存在しない種目",
			modify: func(rec *api.WorkoutRecord) {
				rec.Exercises[1].Exercise.Id = "no-such-exercise"
			},
		},
		{
			name: "終了が開始より前",
			modify: func(rec *api.WorkoutRecord) {
				rec.EndedAt = rec.StartedAt.Add(-time.Minute)
			},
		},
		{
			name: "同じセット ID を2回使う",
			modify: func(rec *api.WorkoutRecord) {
				rec.Exercises[0].Sets[1].Id = rec.Exercises[0].Sets[0].Id
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepos(t)
			rec := newRecord(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC))
			tt.modify(&rec)

			_, err := repo.Put(context.Background(), alice, rec)
			if !errors.Is(err, ErrInvalidWorkout) {
				t.Fatalf("エラー = %v, 期待 ErrInvalidWorkout", err)
			}
			// トランザクションが取り消されて、何も残っていないこと
			if n := countRows(t, repo.pool, "workouts"); n != 0 {
				t.Errorf("workouts の行数 = %d, 期待 0", n)
			}
		})
	}
}

// 5. 新しい順に並ぶ／from・to で絞り込める／ほかの人の記録は出ない
func TestWorkoutRepository_ListOrderAndFilter(t *testing.T) {
	repo := setupRepos(t)
	ctx := context.Background()

	nov1 := newRecord(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC))
	nov2 := newRecord(time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC))
	nov3 := newRecord(time.Date(2026, 11, 3, 10, 0, 0, 0, time.UTC))
	bobs := newRecord(time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC))

	// わざとバラバラの順で入れる
	for _, rec := range []api.WorkoutRecord{nov2, nov3, nov1} {
		if _, err := repo.Put(ctx, alice, rec); err != nil {
			t.Fatalf("Put でエラー: %v", err)
		}
	}
	if _, err := repo.Put(ctx, bob, bobs); err != nil {
		t.Fatalf("bob の Put でエラー: %v", err)
	}

	// 絞り込みなし：新しい順に 11/3, 11/2, 11/1。bob のは入らない
	got, err := repo.List(ctx, alice, nil, nil)
	if err != nil {
		t.Fatalf("List でエラー: %v", err)
	}
	want := []uuid.UUID{nov3.Id, nov2.Id, nov1.Id}
	if len(got) != len(want) {
		t.Fatalf("記録の数 = %d, 期待 %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Id != want[i] {
			t.Errorf("%d 番目の記録が違います（新しい順になっていない）", i)
		}
	}

	// from は「以上」、to は「未満」。11/2 10:00 〜 11/3 10:00 なら 11/2 だけ
	from := nov2.StartedAt
	to := nov3.StartedAt
	got, err = repo.List(ctx, alice, &from, &to)
	if err != nil {
		t.Fatalf("List(from, to) でエラー: %v", err)
	}
	if len(got) != 1 || got[0].Id != nov2.Id {
		t.Errorf("from/to で絞った結果が違います: %d 件", len(got))
	}
}

// 6. 削除すると種目・セットも消える（カスケード）。ほかの人の記録は消せない
func TestWorkoutRepository_Delete(t *testing.T) {
	repo := setupRepos(t)
	ctx := context.Background()
	rec := newRecord(time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC))

	if _, err := repo.Put(ctx, alice, rec); err != nil {
		t.Fatalf("Put でエラー: %v", err)
	}

	// bob が alice の記録を消そうとしても、エラーにはならず何も消えない
	if err := repo.Delete(ctx, bob, rec.Id); err != nil {
		t.Fatalf("bob の Delete でエラー: %v", err)
	}
	if n := countRows(t, repo.pool, "workouts"); n != 1 {
		t.Fatalf("bob の Delete で alice の記録が消えました（workouts = %d）", n)
	}

	// alice が消すと、子と孫も消える
	if err := repo.Delete(ctx, alice, rec.Id); err != nil {
		t.Fatalf("alice の Delete でエラー: %v", err)
	}
	for _, table := range []string{"workouts", "workout_exercises", "workout_sets"} {
		if n := countRows(t, repo.pool, table); n != 0 {
			t.Errorf("%s の行数 = %d, 期待 0", table, n)
		}
	}

	// もう無い記録を消しても、エラーにならない（DELETE は何度送ってもよい）
	if err := repo.Delete(ctx, alice, rec.Id); err != nil {
		t.Errorf("2回目の Delete でエラー: %v", err)
	}
}
