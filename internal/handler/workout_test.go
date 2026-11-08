package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aokitoshihumi/fitness-app-backend/internal/api"
	"github.com/aokitoshihumi/fitness-app-backend/internal/auth"
	"github.com/aokitoshihumi/fitness-app-backend/internal/repository"
	"github.com/google/uuid"
)

// fakeWorkouts は DB の代わり。決めておいた値を返し、受け取った記録を覚えておく
type fakeWorkouts struct {
	created bool
	err     error
	records []api.WorkoutRecord // List が返す値

	called  bool
	gotUID  string
	gotRec  api.WorkoutRecord
	gotFrom *time.Time
	gotTo   *time.Time
	gotID   uuid.UUID
}

func (f *fakeWorkouts) Put(ctx context.Context, firebaseUID string, rec api.WorkoutRecord) (bool, error) {
	f.called = true
	f.gotUID = firebaseUID
	f.gotRec = rec
	return f.created, f.err
}

func (f *fakeWorkouts) List(ctx context.Context, firebaseUID string, from, to *time.Time) ([]api.WorkoutRecord, error) {
	f.called = true
	f.gotUID = firebaseUID
	f.gotFrom, f.gotTo = from, to
	return f.records, f.err
}

func (f *fakeWorkouts) Delete(ctx context.Context, firebaseUID string, id uuid.UUID) error {
	f.called = true
	f.gotUID = firebaseUID
	f.gotID = id
	return f.err
}

const testWorkoutID = "3f2504e0-4f89-41d3-9a0c-0305e82c3301"

const validBody = `{
  "id": "3F2504E0-4F89-41D3-9A0C-0305E82C3301",
  "started_at": "2026-10-03T19:02:00+09:00",
  "ended_at": "2026-10-03T20:10:00+09:00",
  "exercises": [{
    "id": "A1B2C3D4-E5F6-4789-8ABC-DEF012345678",
    "exercise": {"id": "bench-press", "name": "ベンチプレス", "body_part": "chest"},
    "sets": [{"id": "11111111-2222-4333-8444-555555555555", "weight": null, "reps": 10}]
  }]
}`

// 本物の mux（パスから {workoutId} を取り出す役）→ ミドルウェア → PutRecord の順に通す
func callPutRecord(workouts WorkoutStore, id, body string) *httptest.ResponseRecorder {
	return callPutRecordWithUsers(&fakeUsers{}, workouts, id, body)
}

// ユーザー作成の動きも確かめたいテスト用。fakeUsers を外から渡せる
func callPutRecordWithUsers(users UserCreator, workouts WorkoutStore, id, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.Handle("PUT /api/v1/workout-records/{workoutId}",
		auth.Middleware(fakeVerifier{})(http.HandlerFunc(NewWorkoutsHandler(users, workouts).PutRecord)))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/workout-records/"+id, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer any-token")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func callListRecords(workouts WorkoutStore, query string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/workout-records",
		auth.Middleware(fakeVerifier{})(http.HandlerFunc(NewWorkoutsHandler(&fakeUsers{}, workouts).ListRecords)))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/workout-records"+query, nil)
	req.Header.Set("Authorization", "Bearer any-token")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestListRecords_記録がないと空配列(t *testing.T) {
	rec := callListRecords(&fakeWorkouts{records: nil}, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `"workout_records":[]`) {
		t.Errorf("空配列になっていません: %s", rec.Body.String())
	}
}

func TestListRecords_fromとtoを渡す(t *testing.T) {
	workouts := &fakeWorkouts{}
	q := "?from=" + url.QueryEscape("2026-10-01T00:00:00+09:00") +
		"&to=" + url.QueryEscape("2026-11-01T00:00:00+09:00")

	rec := callListRecords(workouts, q)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
	}
	if workouts.gotUID != "user-123" {
		t.Errorf("UID = %q, want %q", workouts.gotUID, "user-123")
	}
	want := time.Date(2026, 10, 1, 0, 0, 0, 0, time.FixedZone("", 9*60*60))
	if workouts.gotFrom == nil || !workouts.gotFrom.Equal(want) {
		t.Errorf("from = %v, want %v", workouts.gotFrom, want)
	}
	if workouts.gotTo == nil {
		t.Error("to が渡っていません")
	}
}

func TestListRecords_日時のいろいろな書き方を読める(t *testing.T) {
	jst := time.FixedZone("", 9*60*60)
	want := time.Date(2026, 10, 1, 0, 0, 0, 0, jst)

	tests := []struct {
		name  string
		query string
	}{
		// iOS の URLQueryItem は + をエンコードしないので、そのまま届く
		{"+をそのまま送る", "?from=2026-10-01T00:00:00+09:00"},
		{"+を%2Bにして送る", "?from=2026-10-01T00:00:00%2B09:00"},
		{"UTC(Z)で送る", "?from=2026-09-30T15:00:00Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workouts := &fakeWorkouts{}

			rec := callListRecords(workouts, tt.query)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
			}
			if workouts.gotFrom == nil || !workouts.gotFrom.Equal(want) {
				t.Errorf("from = %v, want %v", workouts.gotFrom, want)
			}
		})
	}
}

func TestListRecords_日時の形がおかしいと400(t *testing.T) {
	workouts := &fakeWorkouts{}

	rec := callListRecords(workouts, "?from=2026-10-01")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if workouts.called {
		t.Error("おかしなリクエストなのに List が呼ばれました")
	}
}

func TestListRecords_DBエラーなら500(t *testing.T) {
	rec := callListRecords(&fakeWorkouts{err: errors.New("db is down")}, "")

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestPutRecord_新規なら201_置き換えなら200(t *testing.T) {
	tests := []struct {
		name       string
		created    bool
		wantStatus int
	}{
		{"新しく作った", true, http.StatusCreated},
		{"置き換えた", false, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workouts := &fakeWorkouts{created: tt.created}

			rec := callPutRecord(workouts, testWorkoutID, validBody)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if workouts.gotUID != "user-123" {
				t.Errorf("Put に渡った UID = %q, want %q", workouts.gotUID, "user-123")
			}
			if got := workouts.gotRec.Exercises[0].Exercise.Id; got != "bench-press" {
				t.Errorf("Put に渡った種目 ID = %q, want %q", got, "bench-press")
			}
			// 省略した memo は空文字で返る
			if !strings.Contains(rec.Body.String(), `"memo":""`) {
				t.Errorf("body に memo の空文字がありません: %s", rec.Body.String())
			}
		})
	}
}

func TestPutRecord_リクエストがおかしいと400(t *testing.T) {
	tests := []struct {
		name string
		id   string
		body string
	}{
		{"パスが UUID でない", "abc", validBody},
		{"JSON が壊れている", testWorkoutID, `{`},
		{"パスと本文の id が違う", "00000000-0000-4000-8000-000000000000", validBody},
		{"終了が開始より前", testWorkoutID,
			strings.Replace(validBody, "20:10:00", "18:00:00", 1)},
		{"部位が一覧にない", testWorkoutID,
			strings.Replace(validBody, `"chest"`, `"neck"`, 1)},
		{"種目が空", testWorkoutID,
			`{"id":"` + testWorkoutID + `","started_at":"2026-10-03T19:02:00+09:00","ended_at":"2026-10-03T20:10:00+09:00","exercises":[]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workouts := &fakeWorkouts{}

			rec := callPutRecord(workouts, tt.id, tt.body)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
			if workouts.called {
				t.Error("おかしなリクエストなのに Put が呼ばれました")
			}
		})
	}
}

func TestPutRecord_保存のエラーに合わせたステータスを返す(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"ほかの人の ID", repository.ErrWorkoutIDConflict, http.StatusConflict},
		{"DB の制約違反", repository.ErrInvalidWorkout, http.StatusBadRequest},
		{"DB が落ちている", errors.New("db is down"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workouts := &fakeWorkouts{err: tt.err}

			rec := callPutRecord(workouts, testWorkoutID, validBody)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestPutRecord_保存の前にユーザーを作る(t *testing.T) {
	users := &fakeUsers{}
	workouts := &fakeWorkouts{created: true}

	rec := callPutRecordWithUsers(users, workouts, testWorkoutID, validBody)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if users.gotUID != "user-123" {
		t.Errorf("CreateIfNotExists に渡った UID = %q, want %q", users.gotUID, "user-123")
	}
}

func TestPutRecord_ユーザー作成に失敗したら500で保存しない(t *testing.T) {
	users := &fakeUsers{err: errors.New("db is down")}
	workouts := &fakeWorkouts{}

	rec := callPutRecordWithUsers(users, workouts, testWorkoutID, validBody)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if workouts.called {
		t.Error("ユーザーを作れなかったのに Put が呼ばれました")
	}
}

func TestPutRecord_リクエストがおかしいとユーザーを作らない(t *testing.T) {
	users := &fakeUsers{}

	rec := callPutRecordWithUsers(users, &fakeWorkouts{}, testWorkoutID, `{`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if users.gotUID != "" {
		t.Error("おかしなリクエストなのに CreateIfNotExists が呼ばれました")
	}
}

func callDeleteRecords(workouts WorkoutStore, id string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.Handle("DELETE /api/v1/workout-records/{workoutId}",
		auth.Middleware(fakeVerifier{})(http.HandlerFunc(NewWorkoutsHandler(&fakeUsers{}, workouts).DeleteRecords)))

	// DELETE は本文がないので nil を渡す
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/workout-records/"+id, nil)
	req.Header.Set("Authorization", "Bearer any-token")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestDeleteRecords_成功なら204(t *testing.T) {
	fake := &fakeWorkouts{}
	rec := callDeleteRecords(fake, testWorkoutID)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	// 204 は「中身なし」なので、本文が空であることも確かめる
	if rec.Body.Len() != 0 {
		t.Errorf("204 なのに本文があります: %s", rec.Body.String())
	}
	if !fake.called {
		t.Fatal("Delete が呼ばれていません")
	}
	// パスの ID がそのまま repository に渡っているか
	if fake.gotID.String() != testWorkoutID {
		t.Errorf("id = %s, want %s", fake.gotID, testWorkoutID)
	}
}

func TestDeleteRecords_IDがUUIDでないと400(t *testing.T) {
	fake := &fakeWorkouts{}
	rec := callDeleteRecords(fake, "not-a-uuid")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	// DB まで行かずに止まっていること
	if fake.called {
		t.Error("ID がおかしいのに Delete が呼ばれました")
	}
}

func TestDeleteRecords_DBエラーなら500(t *testing.T) {
	rec := callDeleteRecords(&fakeWorkouts{err: errors.New("db down")}, testWorkoutID)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
