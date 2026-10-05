import { OpenAPIHono } from "@hono/zod-openapi";
import { stub } from "./stub";

export const provocations = new OpenAPIHono();
const tag = "煽り";

stub(provocations, { method: "post", path: "/api/v1/provocations", summary: "発動", tag });
stub(provocations, { method: "get", path: "/api/v1/provocations/received", summary: "自分宛の煽り一覧", tag });
stub(provocations, { method: "get", path: "/api/v1/provocations/{id}", summary: "詳細（センキュー一覧を含む）", tag });
stub(provocations, { method: "post", path: "/api/v1/provocations/{id}/thanks", summary: "センキュー送信", tag });
