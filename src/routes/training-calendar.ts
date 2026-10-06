import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { and, asc, eq, gte, lt } from "drizzle-orm";
import type { createDb } from "../db/client";
import { exercises, users, workoutExercises, workoutSessions, workoutSets } from "../db/schema";
import { auth, type AuthVariables } from "../middleware/auth";

type Variables = AuthVariables & { db: ReturnType<typeof createDb> };

export const trainingCalendar = new OpenAPIHono<{ Variables: Variables }>();

const errorBody = (code: string, message: string) => ({ error: { code, message } });

// 仕様書の定義は openapi.yml 側が優先されるので、ここのスキーマは入力チェックが目的。
const route = createRoute({
  method: "get",
  path: "/api/v1/me/training-calendar",
  summary: "自分の1ヶ月分のトレーニング履歴を取得する",
  tags: ["training_calendar"],
  middleware: [auth] as const,
  request: {
    query: z.object({ month: z.string().regex(/^\d{4}-(0[1-9]|1[0-2])$/) }),
  },
  responses: {
    200: { description: "取得成功" },
    400: { description: "month が未指定、または形式が不正" },
    401: { description: "認証失敗" },
  },
});

// 日付 → ユーザーのタイムゾーンでの { 日付, UTCからのオフセット付き日時 } に変換する
const localParts = (date: Date, timeZone: string) => {
  const parts = Object.fromEntries(
    new Intl.DateTimeFormat("en-US", {
      timeZone,
      hourCycle: "h23",
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      timeZoneName: "longOffset",
    })
      .formatToParts(date)
      .map((p) => [p.type, p.value]),
  );
  // longOffset は "GMT+09:00" / "GMT"（UTCのとき）
  const offset = parts.timeZoneName === "GMT" ? "+00:00" : parts.timeZoneName.replace("GMT", "");
  const day = `${parts.year}-${parts.month}-${parts.day}`;
  return { day, iso: `${day}T${parts.hour}:${parts.minute}:${parts.second}${offset}` };
};

trainingCalendar.openapi(
  route,
  async (c) => {
    const { month } = c.req.valid("query");
    const db = c.get("db");
    const userId = c.get("userId");

    const user = await db.select({ timezone: users.timezone }).from(users).where(eq(users.id, userId)).get();
    if (!user) return c.json(errorBody("unauthorized", "認証に失敗しました"), 401);

    // ユーザーのタイムゾーンでの月初〜翌月初は UTC に対して最大±14時間ずれるので、
    // 前後1日ぶん広めにDBから取り、日付はあとで正確に判定して絞る。
    const [year, mon] = month.split("-").map(Number);
    const from = new Date(Date.UTC(year, mon - 1, 1) - 24 * 3600 * 1000).toISOString();
    const to = new Date(Date.UTC(year, mon, 1) + 24 * 3600 * 1000).toISOString();

    const rows = await db
      .select({
        sessionId: workoutSessions.id,
        startedAt: workoutSessions.startedAt,
        workoutExerciseId: workoutExercises.id,
        exerciseId: exercises.id,
        name: exercises.name,
        bodyPart: exercises.bodyPart,
        weightKg: workoutSets.weightKg,
        reps: workoutSets.reps,
      })
      .from(workoutSessions)
      .innerJoin(workoutExercises, eq(workoutExercises.sessionId, workoutSessions.id))
      .innerJoin(exercises, eq(exercises.id, workoutExercises.exerciseId))
      .innerJoin(workoutSets, eq(workoutSets.workoutExerciseId, workoutExercises.id))
      .where(
        and(
          eq(workoutSessions.userId, userId),
          gte(workoutSessions.startedAt, from),
          lt(workoutSessions.startedAt, to),
        ),
      )
      .orderBy(
        asc(workoutSessions.startedAt),
        asc(workoutSessions.id),
        asc(workoutExercises.position),
        asc(workoutSets.position),
      )
      .all();

    type Exercise = {
      exercise_id: number;
      name: string;
      body_part: string;
      sets: { weight_kg: number; reps: number }[];
    };
    type Session = { session_id: number; started_at: string; exercises: Exercise[] };
    type Day = { date: string; body_parts: string[]; sessions: Session[] };

    // 行は「開始時刻 → 種目順 → セット順」に並んでいるので、順に畳み込む。
    const days = new Map<string, Day>();
    const sessions = new Map<number, Session>();
    const exerciseRecords = new Map<number, Exercise>();

    for (const r of rows) {
      let session = sessions.get(r.sessionId);
      if (!session) {
        const local = localParts(new Date(r.startedAt), user.timezone);
        if (!local.day.startsWith(`${month}-`)) continue;
        session = { session_id: r.sessionId, started_at: local.iso, exercises: [] };
        sessions.set(r.sessionId, session);
        let day = days.get(local.day);
        if (!day) {
          day = { date: local.day, body_parts: [], sessions: [] };
          days.set(local.day, day);
        }
        day.sessions.push(session);
      }
      let exercise = exerciseRecords.get(r.workoutExerciseId);
      if (!exercise) {
        exercise = { exercise_id: r.exerciseId, name: r.name, body_part: r.bodyPart, sets: [] };
        exerciseRecords.set(r.workoutExerciseId, exercise);
        session.exercises.push(exercise);
      }
      exercise.sets.push({ weight_kg: r.weightKg, reps: r.reps });
    }

    const result = [...days.values()].sort((a, b) => a.date.localeCompare(b.date));
    for (const day of result) {
      day.body_parts = [...new Set(day.sessions.flatMap((s) => s.exercises.map((e) => e.body_part)))];
    }

    return c.json({ month, days: result }, 200);
  },
  // 入力チェック失敗時は openapi.yml の ErrorResponse の形で 400 を返す
  (result, c) => {
    if (!result.success) {
      return c.json(errorBody("invalid_month", "month は YYYY-MM 形式で指定してください"), 400);
    }
  },
);
