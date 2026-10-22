-- 1回のトレーニング。id はアプリが作った UUID をそのまま使う（誰の何月何日のセッションか）
CREATE TABLE workouts (
    id         UUID        PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- 親が死んだら子も死ぬことにする
    started_at TIMESTAMPTZ NOT NULL,
    ended_at   TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (ended_at >= started_at) 
);
-- idはプライマリーキーなのでインデックスがはられている。実際に誰がどのワークアウトを始めたか絞り込む時に高速化させる
CREATE INDEX workouts_user_id_started_at_idx ON workouts (user_id, started_at DESC);

-- トレーニング内の1種目（ベンチプレスとか）
CREATE TABLE workout_exercises (
    id            UUID    PRIMARY KEY,
    workout_id    UUID    NOT NULL REFERENCES workouts (id) ON DELETE CASCADE,
    -- アプリで並べた順（0 から）
    position      INTEGER NOT NULL CHECK (position >= 0), -- フロントでnセット目がmkgという対応付が必要なためposition(セット目)を作る
    exercise_id   TEXT    NOT NULL REFERENCES exercises (id),
    -- 記録したときの種目名と部位をそのまま残す。あとで種目名が変わっても過去の記録は変わらない
    exercise_name TEXT    NOT NULL,
    body_part     TEXT    NOT NULL REFERENCES body_parts (code),
    UNIQUE (workout_id, position) -- work_outidかつnセット目の組み合わせのUNIQUEは許さない
);

-- 1セット（iOS の WorkoutSet）
CREATE TABLE workout_sets (
    id                  UUID             PRIMARY KEY,
    workout_exercise_id UUID             NOT NULL REFERENCES workout_exercises (id) ON DELETE CASCADE,
    -- 実施順（0 から）
    position            INTEGER          NOT NULL CHECK (position >= 0),
    -- 重量（kg）。自重などで入力していないときは NULL。0kg とは区別する
    weight              DOUBLE PRECISION CHECK (weight >= 0),
    reps                INTEGER          NOT NULL CHECK (reps >= 0),
    -- セットごとのメモ（「補助あり」など）。書いていないときは空文字
    memo                TEXT             NOT NULL DEFAULT '',
    UNIQUE (workout_exercise_id, position)
);
