// Package api は openapi.yml から生成したコードを置く。
// openapi.yml を変更したら `go generate ./...` を実行して api.gen.go を作り直す。
package api

//go:generate go tool oapi-codegen -config config.yaml ../../openapi.yml
