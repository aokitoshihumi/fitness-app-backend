CREATE TABLE exercises (
    id             TEXT PRIMARY KEY,
    name           TEXT NOT NULL UNIQUE,
    body_part_code TEXT NOT NULL REFERENCES body_parts (code)
);
-- body_part_codeは外部キーなのでインデックスをはる
CREATE INDEX exercises_body_part_code_idx ON exercises (body_part_code);
