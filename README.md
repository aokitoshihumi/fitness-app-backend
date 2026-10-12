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
