.PHONY: dev build sqlc

sqlc:
	sqlc generate

dev:
	cd web && npm run dev &
	go run ./cmd/server

build:
	cd web && npm run build
	go build -o server ./cmd/server
