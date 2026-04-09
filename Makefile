.PHONY: dev build sqlc kill

sqlc:
	sqlc generate

dev:
	cd web && npm run dev &
	export $$(grep -v '^#' .env | xargs) && go run ./cmd/server

build:
	cd web && npm run build
	go build -o server ./cmd/server

kill:
	-pkill -f "web/node_modules/.bin/vite"
	-pkill -f "go run.*cmd/server"
	-lsof -ti :8080 | xargs kill -9
	@echo "killed dev processes"
