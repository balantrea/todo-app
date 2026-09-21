.PHONY: up down restart logs build rebuild migrate-up migrate-down clean

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up -d

logs:
	docker compose logs -f

build:
	docker compose up -d --build

rebuild:
	docker compose down
	docker compose up -d --build

migrate-up:
	docker compose up migrate

migrate-down:
	docker compose run --rm migrate \
		-path /schema \
		-database "postgres://$${DB_USER}:$${DB_PASSWORD}@postgres:5432/$${DB_NAME}?sslmode=disable" \
		down 1

clean:
	docker compose down -v