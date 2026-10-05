import { createRoute, z, type OpenAPIHono } from "@hono/zod-openapi";

type StubOptions = {
  method: "get" | "post" | "put" | "patch" | "delete";
  path: string;
  summary: string;
  // Swagger UI 上のグループ名
  tag: string;
  // クエリパラメータ名の一覧（すべて任意の文字列として扱う）
  query?: string[];
};

export const stub = (app: OpenAPIHono, option: StubOptions) => {
  const pathParams = [...option.path.matchAll(/\{(\w+)\}/g)].map((m) => m[1]);

  const route = createRoute({
    method: option.method,
    path: option.path,
    summary: option.summary,
    tags: [option.tag],
    request: {
      params: z.object(
        Object.fromEntries(pathParams.map((p) => [p, z.string()])),
      ),
      query: z.object(
        Object.fromEntries(
          (option.query ?? []).map((q) => [q, z.string().optional()]),
        ),
      ),
    },
    responses: {
      200: { description: "成功（レスポンスのスキーマは未定義）" },
      501: {
        description: "未実装",
        content: {
          "application/json": { schema: z.object({ message: z.string() }) },
        },
      },
    },
  });

  app.openapi(route, (c) => c.json({ message: "Not Implemented" }, 501));
};
