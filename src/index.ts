// Cloudflare Workers はここで書いた `app` を受け取って、リクエストを処理する。
import { OpenAPIHono, createRoute, z } from "@hono/zod-openapi";
import { swaggerUI } from "@hono/swagger-ui";
import { parse } from "yaml";
import openapiYml from "../openapi.yml";
import { createDb } from "./db/client";
import { users } from "./routes/users";
import { friends } from "./routes/friends";
import { gyms } from "./routes/gyms";
import { provocations } from "./routes/provocations";
import { workoutRecords } from "./routes/workout-records";
import { misc } from "./routes/misc";

type Bindings = {
  fitness_app_dev: D1Database;
};

// ミドルウェアで c.set("db", ...) した値の型です。
// 各ルートで c.get("db") と書くと、Drizzle の DB クライアントが取り出せます。
type Variables = {
  db: ReturnType<typeof createDb>;
};

// Hono アプリ本体。OpenAPIHono は Hono に「API仕様書の自動生成」機能を足したものです。
// Bindings と Variables の型を教えておくと、補完やエラーチェックが効きます。
const app = new OpenAPIHono<{ Bindings: Bindings; Variables: Variables }>();

// ミドルウェア: すべてのリクエスト（"*"）の前に実行される処理です。
// ここで DB クライアントを作り、各ルートから使えるように c.set で保存しています。
app.use("*", async (c, next) => {
  c.set("db", createDb(c.env.fitness_app_dev));
  // next() で次の処理（ルートの処理）に進みます。忘れると先に進まないので注意。
  await next();
});

// GET / にアクセスされたときの処理です。動作確認用の「Hello Hono!」を返します。
// createRoute で「パス・メソッド・返すもの」を定義すると、そのまま仕様書にも載ります。
// 新しいエンドポイントも、この書き方で追加していきます。
const rootRoute = createRoute({
  method: "get",
  path: "/",
  summary: "動作確認",
  responses: {
    200: {
      description: "Hello Hono!",
      content: { "text/plain": { schema: z.string() } },
    },
  },
});

app.openapi(rootRoute, (c) => {
  return c.text("Hello Hono!");
});

// docs/API.md のエンドポイント。まだ中身は未実装（501を返す）で、仕様書にだけ載っています。
app.route("/", users);
app.route("/", friends);
app.route("/", gyms);
app.route("/", provocations);
app.route("/", workoutRecords);
app.route("/", misc);

// 仕様書(JSON)。/doc で取得できます。
// コードから生成した仕様（未実装のスタブ）に、openapi.yml で固めた定義を上書きして返します。
// ymlのパスは servers の /api/v1 からの相対なので、ここで /api/v1 を付けて揃えます。
const API_PREFIX = "/api/v1";
const yml = parse(openapiYml) as {
  paths: Record<string, unknown>;
  components?: Record<string, Record<string, unknown>>;
  tags?: unknown[];
  security?: unknown[];
};

app.get("/doc", (c) => {
  const generated = app.getOpenAPIDocument({
    openapi: "3.0.0",
    info: { title: "Fitness App API", version: "1.0.0" },
  });

  // ymlで定義したパスは、スタブ側の同じパスを置き換える
  const paths: Record<string, unknown> = { ...generated.paths };
  for (const [path, item] of Object.entries(yml.paths)) {
    paths[API_PREFIX + path] = item;
  }

  const components: Record<string, Record<string, unknown>> = { ...generated.components };
  for (const [kind, defs] of Object.entries(yml.components ?? {})) {
    components[kind] = { ...components[kind], ...defs };
  }

  return c.json({
    ...generated,
    paths,
    components,
    tags: [...(yml.tags ?? []), ...(generated.tags ?? [])],
    security: yml.security,
  });
});

// Swagger UI。ブラウザで /ui を開くと仕様書を見たり試したりできます。
app.get("/ui", swaggerUI({ url: "/doc" }));

// Cloudflare Workers が実行するためにアプリを export しています。
export default app;
