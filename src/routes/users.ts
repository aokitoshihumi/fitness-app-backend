import { OpenAPIHono } from "@hono/zod-openapi";
import { stub } from "./stub";

export const users = new OpenAPIHono();
const tag = "ユーザー";

stub(users, { method: "put", path: "/api/v1/users/me", summary: "初回登録（冪等）", tag });
stub(users, { method: "get", path: "/api/v1/users/me", summary: "自分のプロフィール", tag });
stub(users, { method: "patch", path: "/api/v1/users/me", summary: "プロフィール・公開設定の更新", tag });
stub(users, { method: "delete", path: "/api/v1/users/me", summary: "退会", tag });
stub(users, { method: "get", path: "/api/v1/users/{id}", summary: "他ユーザーのプロフィール", tag });
stub(users, { method: "get", path: "/api/v1/users/{id}/contributions", summary: "草データ", tag });
stub(users, { method: "get", path: "/api/v1/users/{id}/deviation-scores", summary: "部位別偏差値", tag });
