import { sql } from "drizzle-orm";
import { index, integer, real, sqliteTable, text } from "drizzle-orm/sqlite-core";

// users テーブル
export const users = sqliteTable("users", {
  // id: 主キー
  id: integer("id").primaryKey({ autoIncrement: true }),
  // name: ユーザー名。非NULL
  name: text("name").notNull(),
  // timezone: IANAタイムゾーン名。「何月何日か」の判定に使う
  timezone: text("timezone").notNull().default("Asia/Tokyo"),
  // created_at: 作成日時
  createdAt: text("created_at")
    .notNull()
    .default(sql`(CURRENT_TIMESTAMP)`),
});

// exercises テーブル: 種目マスタ
export const exercises = sqliteTable("exercises", {
  id: integer("id").primaryKey({ autoIncrement: true }),
  // name: 種目名（例: ベンチプレス）
  name: text("name").notNull(),
  // body_part: 部位コード（chest / back / shoulders / arms / legs / abs）
  bodyPart: text("body_part", {
    enum: ["chest", "back", "shoulders", "arms", "legs", "abs"],
  }).notNull(),
});

// workout_sessions テーブル: 1回のトレーニング
export const workoutSessions = sqliteTable(
  "workout_sessions",
  {
    id: integer("id").primaryKey({ autoIncrement: true }),
    userId: integer("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),
    // started_at: 開始日時。UTCのISO 8601（例: 2026-10-03T01:15:00.000Z）で保存する
    startedAt: text("started_at").notNull(),
  },
  (t) => [index("workout_sessions_user_started_idx").on(t.userId, t.startedAt)],
);

// workout_exercises テーブル: セッション内で実施した種目
export const workoutExercises = sqliteTable(
  "workout_exercises",
  {
    id: integer("id").primaryKey({ autoIncrement: true }),
    sessionId: integer("session_id")
      .notNull()
      .references(() => workoutSessions.id, { onDelete: "cascade" }),
    exerciseId: integer("exercise_id")
      .notNull()
      .references(() => exercises.id),
    // position: セッション内の実施順（0始まり）
    position: integer("position").notNull(),
  },
  (t) => [index("workout_exercises_session_idx").on(t.sessionId)],
);

// workout_sets テーブル: 種目ごとのセット記録
export const workoutSets = sqliteTable(
  "workout_sets",
  {
    id: integer("id").primaryKey({ autoIncrement: true }),
    workoutExerciseId: integer("workout_exercise_id")
      .notNull()
      .references(() => workoutExercises.id, { onDelete: "cascade" }),
    // position: 種目内の実施順（0始まり）
    position: integer("position").notNull(),
    // weight_kg: 重量。自重種目は 0
    weightKg: real("weight_kg").notNull(),
    reps: integer("reps").notNull(),
  },
  (t) => [index("workout_sets_workout_exercise_idx").on(t.workoutExerciseId)],
);
