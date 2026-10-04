import { sql } from "drizzle-orm";
import { integer, sqliteTable, text } from "drizzle-orm/sqlite-core";

// users テーブル
export const users = sqliteTable("users", {
  // id: 主キー
  id: integer("id").primaryKey({ autoIncrement: true }),
  // name: ユーザー名。非NULL
  name: text("name").notNull(),
  // created_at: 作成日時
  createdAt: text("created_at")
    .notNull()
    .default(sql`(CURRENT_TIMESTAMP)`),
});
