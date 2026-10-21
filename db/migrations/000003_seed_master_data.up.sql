-- iOS アプリ（Exercise.swift）の BodyPart と Exercise.presets と同じ値にする
INSERT INTO body_parts (code, label) VALUES
    ('chest',     '胸'),
    ('back',      '背中'),
    ('legs',      '脚'),
    ('shoulders', '肩'),
    ('arms',      '腕'),
    ('abs',       '腹筋');

INSERT INTO exercises (id, name, body_part_code) VALUES
    ('bench-press',            'ベンチプレス',               'chest'),
    ('incline-dumbbell-press', 'インクラインダンベルプレス', 'chest'),
    ('cable-fly',              'ケーブルフライ',             'chest'),
    ('deadlift',               'デッドリフト',               'back'),
    ('lat-pulldown',           'ラットプルダウン',           'back'),
    ('bent-over-row',          'ベントオーバーロウ',         'back'),
    ('squat',                  'スクワット',                 'legs'),
    ('leg-press',              'レッグプレス',               'legs'),
    ('leg-curl',               'レッグカール',               'legs'),
    ('shoulder-press',         'ショルダープレス',           'shoulders'),
    ('side-raise',             'サイドレイズ',               'shoulders'),
    ('dumbbell-curl',          'ダンベルカール',             'arms'),
    ('triceps-pushdown',       'トライセプスプッシュダウン', 'arms'),
    ('crunch',                 'クランチ',                   'abs'),
    ('leg-raise',              'レッグレイズ',               'abs');
