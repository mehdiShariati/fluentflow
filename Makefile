.PHONY: up down logs build test verify docs

# Full stack: Postgres + API + Web + LiveKit + agent (requires Docker)
up:
	docker compose up --build

# Go unit tests (exclude web/node_modules nested Go packages)
test:
	go test -buildvcs=false ./cmd/... ./internal/...

# Build MkDocs site to ./site (requires: pip install -r requirements-docs.txt)
docs:
	mkdocs build

# Go vet + tests + Next.js lint/build (requires Node/npm for web)
verify:
	go vet -buildvcs=false ./cmd/... ./internal/...
	go test -buildvcs=false ./cmd/... ./internal/...
	npm --prefix web run lint
	npm --prefix web run build

down:
	docker compose down

logs:
	docker compose logs -f

build:
	docker compose build
