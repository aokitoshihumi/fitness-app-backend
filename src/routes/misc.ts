import { OpenAPIHono } from "@hono/zod-openapi";
import { stub } from "./stub";

export const misc = new OpenAPIHono();

stub(misc, { method: "get", path: "/api/v1/timeline", summary: "フレンドの活動", tag: "タイムライン", query: ["limit", "cursor"] });
stub(misc, { method: "post", path: "/api/v1/body-weight-logs", summary: "記録", tag: "体重" });
stub(misc, { method: "get", path: "/api/v1/body-weight-logs", summary: "推移", tag: "体重" });
stub(misc, { method: "get", path: "/api/v1/exercises", summary: "種目一覧", tag: "マスタ" });
stub(misc, { method: "get", path: "/api/v1/body-parts", summary: "部位一覧", tag: "マスタ" });
