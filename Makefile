# .env があれば読み込み、中の値を環境変数として下のコマンドに渡す
-include .env
export

DB_HOST ?= localhost
DB_PORT ?= 5432
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# golang-migrate はインストールせず、バージョンを固定して go run で動かす
MIGRATE_VERSION := v4.20.1
MIGRATE := go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

.PHONY: db-up db-down migrate-up migrate-down migrate-version migrate-create run

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