CREATE TABLE users (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Firebase Auth のユーザー ID。ID トークンからこの値を取り出して、誰のリクエストかを決める
    firebase_uid TEXT        NOT NULL UNIQUE,
    -- 画面に表示する名前。ユーザーが自由につけるので、他の人と同じでもよい。初回登録時はまだ決まっていないので空を許す
    display_name TEXT,
    -- 他のユーザーから見える ID（@oren など）。検索に使うので重複を許さない。初回登録時はまだ決まっていないので空を許す
    handle       TEXT        UNIQUE,
    time_zone    TEXT        NOT NULL DEFAULT 'Asia/Tokyo',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
