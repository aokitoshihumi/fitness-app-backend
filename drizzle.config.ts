import { defineConfig } from "drizzle-kit";

export default defineConfig({
  // Cloudflare D1はSQLite
  dialect: "sqlite",
  schema: "./src/db/schema.ts",
  // 生成された SQL やメタ情報を置く場所
  out: "./drizzle",
});
