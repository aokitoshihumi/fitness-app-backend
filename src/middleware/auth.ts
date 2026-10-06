import { createMiddleware } from "hono/factory";

export type AuthVariables = {
  userId: number;
};

// TODO: 仮実装。本来は Authorization ヘッダーの Firebase ID トークンを検証して uid からユーザーを特定する。
// Firebase 連携ができるまでは、開発用に X-Debug-User-Id ヘッダーの値をユーザーIDとして扱う。
// ※誰でもなりすませるので、本番に出す前に必ず置き換えること。
export const auth = createMiddleware<{ Variables: AuthVariables }>(async (c, next) => {
  const userId = Number(c.req.header("X-Debug-User-Id"));
  if (!Number.isInteger(userId) || userId <= 0) {
    return c.json({ error: { code: "unauthorized", message: "認証に失敗しました" } }, 401);
  }
  c.set("userId", userId);
  await next();
});
