CREATE TABLE exercises (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name         TEXT   NOT NULL UNIQUE,
    body_part_id BIGINT NOT NULL REFERENCES body_parts (id)
);
-- body_parts_idは外部キーなのでインデックスをはる
CREATE INDEX exercises_body_part_id_idx ON exercises (body_part_id);
