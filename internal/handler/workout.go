package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aokitoshihumi/fitness-app-backend/internal/api"
	"github.com/aokitoshihumi/fitness-app-backend/internal/auth"
	"github.com/aokitoshihumi/fitness-app-backend/internal/repository"
	"github.com/google/uuid"
)

type WorkoutSaver interface {
	Put(ctx context.Context, firebaseUID string, rec api.WorkoutRecord) (created bool, err error)
}

type WorkoutsHandler struct {
	workouts WorkoutSaver
}

func NewWorkoutsHandler(workouts WorkoutSaver) *WorkoutsHandler {
	return &WorkoutsHandler{workouts: workouts}
}

// PUT /api/v1/workout-records/{workoutId} のハンドラ
// auth.Middlewareの後ろに置く
// 新規作成は201、置き換えは200。どちらも保存した記録を返す
func (h *WorkoutsHandler) PutRecord(w http.ResponseWriter, r *http.Request) {
	uid, ok := auth.UIDFromContext(r.Context())
	if !ok {
		log.Println("PUT /workout-records: contextにUIDがありません。auth.Middlewareを通していますか")
		writeInternalError(w)
		return
	}

	// 1. パスの {workoutId} を取り出して UUID として読む
	pathID, err := uuid.Parse(r.PathValue("workoutId"))
	if err != nil {
		writeBadRequest(w)
		return
	}

	// 2. 本文の JSON を api.WorkoutRecord に読み込む
	var rec api.WorkoutRecord
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		writeBadRequest(w)
		return
	}

	// 3. 中身が仕様どおりか確かめる
	if err := validateWorkoutRecord(pathID, rec); err != nil {
		log.Printf("PUT /workout-records: %v", err)
		writeBadRequest(w)
		return
	}

	// 4. 保存する
	created, err := h.workouts.Put(r.Context(), uid, rec)
	switch {
	case errors.Is(err, repository.ErrWorkoutIDConflict):
		writeError(w, http.StatusConflict, "id_conflict", "この ID はすでに使われています")
		return
	case errors.Is(err, repository.ErrInvalidWorkout):
		log.Printf("PUT /workout-records: %v", err)
		writeBadRequest(w)
		return
	case err != nil:
		log.Printf("PUT /workout-records: %v", err)
		writeInternalError(w)
		return
	}

	// 5. 保存した記録を返す。memo を省略していたら、保存した値（空文字）に合わせる
	for i := range rec.Exercises {
		for j := range rec.Exercises[i].Sets {
			if rec.Exercises[i].Sets[j].Memo == nil {
				empty := ""
				rec.Exercises[i].Sets[j].Memo = &empty
			}
		}
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(rec)
}

// validateWorkoutRecord は openapi.yml に書いた決まりを確かめる。
// DB の制約でも防げるものもあるが、ここで先に弾くとエラーの理由がわかりやすい
func validateWorkoutRecord(pathID uuid.UUID, rec api.WorkoutRecord) error {
	if rec.Id != pathID {
		return errors.New("本文の id とパスの workoutId が違います")
	}
	if rec.StartedAt.IsZero() || rec.EndedAt.IsZero() {
		return errors.New("started_at と ended_at は必須です")
	}
	if rec.EndedAt.Before(rec.StartedAt) {
		return errors.New("ended_at が started_at より前です")
	}
	if len(rec.Exercises) == 0 {
		return errors.New("exercises が空です")
	}
	for _, ex := range rec.Exercises {
		if ex.Exercise.Id == "" || ex.Exercise.Name == "" {
			return errors.New("exercise の id と name は必須です")
		}
		if !ex.Exercise.BodyPart.Valid() {
			return errors.New("body_part が正しくありません")
		}
		if len(ex.Sets) == 0 {
			return errors.New("sets が空です")
		}
		for _, set := range ex.Sets {
			if set.Reps < 0 {
				return errors.New("reps がマイナスです")
			}
			if set.Weight != nil && *set.Weight < 0 {
				return errors.New("weight がマイナスです")
			}
		}
	}
	return nil
}

func writeBadRequest(w http.ResponseWriter) {
	writeError(w, http.StatusBadRequest, "invalid_request", "リクエストの形式が正しくありません")
}

// writeError は {"error":{"code":..., "message":...}} の形でエラーを返す
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var body api.ErrorResponse
	body.Error.Code = code
	body.Error.Message = message
	json.NewEncoder(w).Encode(body)
}
