import { OpenAPIHono } from "@hono/zod-openapi";
import { stub } from "./stub";

export const gyms = new OpenAPIHono();
const tag = "ジム・チェックイン";

stub(gyms, { method: "get", path: "/api/v1/gyms", summary: "周辺ジム検索", tag, query: ["lat", "lng"] });
stub(gyms, { method: "post", path: "/api/v1/gyms", summary: "ジム新規登録", tag });
stub(gyms, { method: "post", path: "/api/v1/gyms/{id}/checkins", summary: "チェックイン", tag });
stub(gyms, { method: "post", path: "/api/v1/checkins/{id}/ping", summary: "滞在確認（ハートビート）", tag });
stub(gyms, { method: "post", path: "/api/v1/checkins/{id}/checkout", summary: "チェックアウト", tag });
stub(gyms, { method: "get", path: "/api/v1/checkins/current", summary: "現在のチェックイン（煽り可否を含む）", tag });
