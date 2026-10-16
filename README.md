# fitness-app-backend

筋トレアプリのバックエンド（Go）。

API の仕様は [openapi.yml](openapi.yml) を正とし、そこから Go のコードを生成します。

## コード生成

`openapi.yml` を変更したら、次を実行して `internal/api/api.gen.go` を作り直します。

```txt
go generate ./...
```

`api.gen.go` は自動生成されるので、手で編集しません。

## ビルド確認

```txt
go build ./...
```

## 旧 TypeScript 版

Hono + Cloudflare Workers で書いていた版は、タグ `ts-final` に残っています。

## 環境変数

| 名前          | 必須 | デフォルト  | 説明                                                   |
| ------------- | ---- | ----------- | ------------------------------------------------------ |
| `PORT`        |      | `8080`      | HTTPサーバーのポート（Cloud Run では自動で設定される） |
| `DB_HOST`     |      | `localhost` | DBのホスト                                             |
| `DB_PORT`     |      | `5432`      | DBのポート                                             |
| `DB_USER`     | ✓    |             | DBのユーザー名                                         |
| `DB_PASSWORD` | ✓    |             | DBのパスワード                                         |
| `DB_NAME`     | ✓    |             | DB名                                                   |

ローカルでは `.env` に書きます。本番（Cloud Run）では Cloud Run の環境変数に設定し、
`DB_PASSWORD` は Secret Manager から渡します。
