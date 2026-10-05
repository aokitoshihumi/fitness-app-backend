import { OpenAPIHono } from "@hono/zod-openapi";
import { stub } from "./stub";

export const workoutRecords = new OpenAPIHono();
const tag = "筋トレ記録";

stub(workoutRecords, { method: "post", path: "/api/v1/workout-records", summary: "記録登録（セッション単位で一括）", tag });
stub(workoutRecords, { method: "get", path: "/api/v1/workout-records", summary: "履歴", tag, query: ["exercise_id", "from", "to"] });
stub(workoutRecords, { method: "get", path: "/api/v1/workout-records/{id}", summary: "詳細", tag });
stub(workoutRecords, { method: "patch", path: "/api/v1/workout-records/{id}", summary: "修正", tag });
stub(workoutRecords, { method: "delete", path: "/api/v1/workout-records/{id}", summary: "削除", tag });
