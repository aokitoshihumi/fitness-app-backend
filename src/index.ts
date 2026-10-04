// Cloudflare Workers はここで書いた `app` を受け取って、リクエストを処理する。
import { Hono } from "hono";
import { createDb } from "./db/client";

type Bindings = {
  fitness_app_dev: D1Database;
};

// ミドルウェアで c.set("db", ...) した値の型です。
// 各ルートで c.get("db") と書くと、Drizzle の DB クライアントが取り出せます。
type Variables = {
  db: ReturnType<typeof createDb>;
};

// Hono アプリ本体。Bindings と Variables の型を教えておくと、補完やエラーチェックが効きます。
const app = new Hono<{ Bindings: Bindings; Variables: Variables }>();

// ミドルウェア: すべてのリクエスト（"*"）の前に実行される処理です。
// ここで DB クライアントを作り、各ルートから使えるように c.set で保存しています。
app.use("*", async (c, next) => {
  c.set("db", createDb(c.env.fitness_app_dev));
  // next() で次の処理（ルートの処理）に進みます。忘れると先に進まないので注意。
  await next();
});

// GET / にアクセスされたときの処理です。動作確認用の「Hello Hono!」を返します。
app.get("/", (c) => {
  return c.text("Hello Hono!");
});

// Cloudflare Workers が実行するためにアプリを export しています。
export default app;
