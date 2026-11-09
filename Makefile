# .env があれば読み込み、中の値を環境変数として下のコマンドに渡す
-include .env
export

DB_HOST ?= localhost
DB_PORT ?= 5432
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# テスト用の DB。開発用のデータを消さないように、同じ PostgreSQL の中に別の DB を作る
TEST_DB_NAME ?= $(DB_NAME)_test
TEST_DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(TEST_DB_NAME)?sslmode=disable
# 上の export で全部の変数が環境変数になるが、これだけは外す。
# 外さないと make test でも DB のテストが動いてしまう（test-db でだけ渡す）
unexport TEST_DATABASE_URL

# golang-migrate はインストールせず、バージョンを固定して go run で動かす
MIGRATE_VERSION := v4.20.1
MIGRATE := go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

.PHONY: db-up db-down migrate-up migrate-down migrate-version migrate-create run token test test-db-setup test-db

## ローカルの PostgreSQL を起動する
db-up:
	docker compose up -d --wait

## ローカルの PostgreSQL を止める（データは残る）
db-down:
	docker compose down

## まだ流していないマイグレーションをすべて流す
migrate-up:
	@$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" up

## 最後に流したマイグレーションを1つだけ戻す
migrate-down:
	@$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" down 1

## 今どのマイグレーションまで流したかを表示する
migrate-version:
	@$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" version

## 新しいマイグレーションファイルを作る（例: make migrate-create name=create_users）
migrate-create:
	@$(MIGRATE) create -ext sql -dir db/migrations -seq $(name)

run:
	go run ./cmd/server

## テスト用ユーザーで Firebase にログインし、ID トークン（1時間有効）だけを表示する
## 使い方: export TOKEN=$(make -s token)
## .env に FIREBASE_API_KEY / TEST_USER_EMAIL / TEST_USER_PASSWORD が必要
token:
	@res=$$(jq -n --arg email "$$TEST_USER_EMAIL" --arg password "$$TEST_USER_PASSWORD" \
		'{email: $$email, password: $$password, returnSecureToken: true}' \
	| curl -s -X POST \
		"https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$$FIREBASE_API_KEY" \
		-H "Content-Type: application/json" -d @-); \
	echo "$$res" | jq -er '.idToken // empty' 2>/dev/null \
	|| { echo "トークンを取れませんでした: $$(echo "$$res" | jq -r '.error.message // .')" >&2; exit 1; }

## DB を使わないテストだけ流す（repository のテストはスキップされる）
test:
	go test ./...

## テスト用の DB が無ければ作り、マイグレーションを流す
test-db-setup:
	@docker compose exec -T db psql -U "$(DB_USER)" -d postgres -tAc \
		"SELECT 1 FROM pg_database WHERE datname = '$(TEST_DB_NAME)'" | grep -q 1 \
	|| docker compose exec -T db psql -U "$(DB_USER)" -d postgres -c 'CREATE DATABASE "$(TEST_DB_NAME)"'
	@$(MIGRATE) -path db/migrations -database "$(TEST_DATABASE_URL)" up

## DB を使うテストも含めて、全部のテストを流す（先に make db-up）
test-db: test-db-setup
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -count=1 ./...
