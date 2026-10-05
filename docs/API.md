■ ユーザー
PUT /api/v1/users/me 初回登録（冪等）
GET /api/v1/users/me 自分のプロフィール
PATCH /api/v1/users/me プロフィール・公開設定の更新
DELETE /api/v1/users/me 退会
GET /api/v1/users/{id} 他ユーザーのプロフィール
GET /api/v1/users/{id}/contributions 草データ
GET /api/v1/users/{id}/deviation-scores 部位別偏差値

■ フレンド
GET /api/v1/friendships フレンド一覧（アクティブ状態を含む）
DELETE /api/v1/friendships/{userId} フレンド解除
POST /api/v1/friend-requests 申請
GET /api/v1/friend-requests 届いた申請一覧
POST /api/v1/friend-requests/{id}/accept 承認
POST /api/v1/friend-requests/{id}/reject 拒否

■ ジム・チェックイン
GET /api/v1/gyms?lat=&lng= 周辺ジム検索
POST /api/v1/gyms ジム新規登録
POST /api/v1/gyms/{id}/checkins チェックイン
POST /api/v1/checkins/{id}/ping 滞在確認（ハートビート）
POST /api/v1/checkins/{id}/checkout チェックアウト
GET /api/v1/checkins/current 現在のチェックイン（煽り可否を含む）

■ 煽り
POST /api/v1/provocations 発動
GET /api/v1/provocations/received 自分宛の煽り一覧
GET /api/v1/provocations/{id} 詳細（センキュー一覧を含む）
POST /api/v1/provocations/{id}/thanks センキュー送信

■ 筋トレ記録
POST /api/v1/workout-records 記録登録（セッション単位で一括）
GET /api/v1/workout-records?exercise_id=&from=&to= 履歴
GET /api/v1/workout-records/{id} 詳細
PATCH /api/v1/workout-records/{id} 修正
DELETE /api/v1/workout-records/{id} 削除

■ タイムライン
GET /api/v1/timeline?limit=&cursor= フレンドの活動

■ 体重
POST /api/v1/body-weight-logs 記録
GET /api/v1/body-weight-logs 推移

■ マスタ
GET /api/v1/exercises 種目一覧
GET /api/v1/body-parts 部位一覧
