CREATE TABLE users (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Firebase Auth のユーザー ID。ID トークンからこの値を取り出して、誰のリクエストかを決める
    firebase_uid TEXT        NOT NULL UNIQUE,
    -- 他のユーザーから見える ID（@oren など）。初回登録時はまだ決まっていないので空を許す
    handle_name TEXT UK UNIQUE -- ハンドルネーム
    name  TEXT        UNIQUE, -- 検索用の名前的なもの
    time_zone    TEXT        NOT NULL DEFAULT 'Asia/Tokyo',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
