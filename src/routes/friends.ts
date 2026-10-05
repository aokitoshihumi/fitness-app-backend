import { OpenAPIHono } from "@hono/zod-openapi";
import { stub } from "./stub";

export const friends = new OpenAPIHono();
const tag = "フレンド";

stub(friends, { method: "get", path: "/api/v1/friendships", summary: "フレンド一覧（アクティブ状態を含む）", tag });
stub(friends, { method: "delete", path: "/api/v1/friendships/{userId}", summary: "フレンド解除", tag });
stub(friends, { method: "post", path: "/api/v1/friend-requests", summary: "申請", tag });
stub(friends, { method: "get", path: "/api/v1/friend-requests", summary: "届いた申請一覧", tag });
stub(friends, { method: "post", path: "/api/v1/friend-requests/{id}/accept", summary: "承認", tag });
stub(friends, { method: "post", path: "/api/v1/friend-requests/{id}/reject", summary: "拒否", tag });
