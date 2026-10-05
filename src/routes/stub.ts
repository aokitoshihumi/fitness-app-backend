import { createRoute, z, type OpenAPIHono } from "@hono/zod-openapi";

type StubOptions = {
  method: "get" | "post" | "put" | "patch" | "delete";
  // パス中の {id} などは、パスパラメータとして自動で仕様書に載ります。
  path: string;
  summary: string;
  // Swagger UI 上のグループ名
  tag: string;
  // クエリパラメータ名の一覧（すべて任意の文字列として扱う）
  query?: string[];
};

// 「仕様書には載せるが、中身はまだ実装していない」エンドポイントを登録するための関数です。
// 実装するときは、この stub(...) を createRoute + app.openapi(...) に置き換えて、
// request / responses にスキーマを書いていきます。
export const stub = (app: OpenAPIHono, o: StubOptions) => {
  const pathParams = [...o.path.matchAll(/\{(\w+)\}/g)].map((m) => m[1]);

  const route = createRoute({
    method: o.method,
    path: o.path,
    summary: o.summary,
    tags: [o.tag],
    request: {
      params: z.object(Object.fromEntries(pathParams.map((p) => [p, z.string()]))),
      query: z.object(Object.fromEntries((o.query ?? []).map((q) => [q, z.string().optional()]))),
    },
    responses: {
      200: { description: "成功（レスポンスのスキーマは未定義）" },
      501: {
        description: "未実装",
        content: { "application/json": { schema: z.object({ message: z.string() }) } },
      },
    },
  });

  app.openapi(route, (c) => c.json({ message: "Not Implemented" }, 501));
};
