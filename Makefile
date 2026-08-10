.PHONY: up down backend frontend test lint

up:
	docker compose up -d

down:
	docker compose down

backend:
	cd backend && go run ./cmd/api

frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./...
	cd frontend && npm test -- --run

lint:
	cd backend && golangci-lint run ./... || go vet ./...
	cd frontend && npm run lint
