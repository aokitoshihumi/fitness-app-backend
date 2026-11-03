package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aokitoshihumi/fitness-app-backend/internal/api"
	"github.com/aokitoshihumi/fitness-app-backend/internal/auth"
	"github.com/aokitoshihumi/fitness-app-backend/internal/repository"
)

// fakeWorkouts は DB の代わり。決めておいた値を返し、受け取った記録を覚えておく
type fakeWorkouts struct {
	created bool
	err     error

	called bool
	gotUID string
	gotRec api.WorkoutRecord
}

func (f *fakeWorkouts) Put(ctx context.Context, firebaseUID string, rec api.WorkoutRecord) (bool, error) {
	f.called = true
	f.gotUID = firebaseUID
	f.gotRec = rec
	return f.created, f.err
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
func callPutRecord(workouts WorkoutSaver, id, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.Handle("PUT /api/v1/workout-records/{workoutId}",
		auth.Middleware(fakeVerifier{})(http.HandlerFunc(NewWorkoutsHandler(workouts).PutRecord)))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/workout-records/"+id, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer any-token")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
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
