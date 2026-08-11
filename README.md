# Equipment Availability System

Fleet availability and dispatch for **Clements Contractors** — project-centric equipment scheduling (not rental billing).

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go, Gin, GORM, PostgreSQL |
| Frontend | React, TypeScript, Vite, Tailwind, TanStack Query |
| Auth | JWT access token (memory) + refresh cookie |

## Quick start

### Prerequisites

- Go 1.22+
- Node 20+
- Docker (Postgres)

### Database

```bash
docker compose up -d
```

### Backend

```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

API: `http://localhost:8080`  
Default seed: `admin@clements.local` / `admin@123`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

App: `http://localhost:5173`

## Documentation

See [docs/README.md](docs/README.md) for architecture, data model, API, and algorithms.

## Make targets

```bash
make up          # postgres
make backend     # run API
make frontend    # run Vite
make test        # go test + frontend test
make lint        # golangci-lint + eslint
```

## License

Internal product planning / proprietary.
