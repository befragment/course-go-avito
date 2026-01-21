APP_NAME := courier-service

SERVICE_CMD := ./cmd/service
WORKER_CMD  := ./cmd/worker
POSTGRES_CONTAINER_NAME := my-postgres

DOCKER_COMPOSE := docker compose

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - build both service and worker binaries"
	@echo "  build-service   - build HTTP service binary"
	@echo "  build-worker    - build worker binary"
	@echo "  run-service     - run HTTP service with go run"
	@echo "  run-worker      - run worker with go run"
	@echo "  test            - run go test ./..."
	@echo "  test-nocache    - run go test ./... -count=1"
	@echo "  lint            - run go vet ./..."
	@echo "  migrate-create  - create new migration (usage: make migrate-create NAME=migration_name)"
	@echo "  migrate-up      - apply DB migrations (scripts/migrate.sh up)"
	@echo "  migrate-down    - rollback last DB migration (scripts/migrate.sh down)"
	@echo "  migrate-status  - show DB migrations status"
	@echo "  apply-seed      - apply seed data to DB"
	@echo "  dc-up           - run docker compose up -d"
	@echo "  dc-build        - run docker compose up -d --build"
	@echo "  dc-down         - run docker compose down"
	@echo "  dc-logs         - logs without prometheus and grafana"
	@echo "  dc-psql         - run psql in postgres container"
	
.PHONY: build
build: build-service build-worker

.PHONY: build-service
build-service:
	go build -o bin/service $(SERVICE_CMD)

.PHONY: build-worker
build-worker:
	go build -o bin/worker $(WORKER_CMD)

.PHONY: run-service
run-service:
	go run $(SERVICE_CMD)

.PHONY: run-worker
run-worker:
	go run $(WORKER_CMD)

.PHONY: test
test:
	go test ./...

.PHONY: test-nocache
test-nocache:
	go test ./... -count=1

.PHONY: lint
lint:
	go vet ./...

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	goose -dir ./migrations create $(NAME) sql

.PHONY: migrate-up
migrate-up:
	sh ./scripts/migrate.sh up

.PHONY: migrate-down
migrate-down:
	sh ./scripts/migrate.sh down

.PHONY: migrate-status
migrate-status:
	sh ./scripts/migrate.sh status

.PHONY: apply-seed
apply-seed:
	sh ./scripts/seed.sh

.PHONY: dc-up
dc-up:
	$(DOCKER_COMPOSE) up -d

.PHONY: dc-build
dc-build:
	$(DOCKER_COMPOSE) up -d --build

.PHONY: dc-down
dc-down:
	$(DOCKER_COMPOSE) down

.PHONY: dc-logs
dc-logs:
	$(DOCKER_COMPOSE) logs -f service-courier service-courier-worker

.PHONY: dc-psql
dc-psql:
	sh ./scripts/psql-container.sh