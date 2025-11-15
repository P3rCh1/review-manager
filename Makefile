.PHONY: migrate-up migrate-down migrate-version migrate-force down up

default: up

migrate-up:
	docker compose run --rm postgres_migrate up

migrate-down:
	docker compose run --rm postgres_migrate down

migrate-version:
	docker compose run --rm postgres_migrate version

migrate-force:
	docker compose run --rm postgres_migrate force $(VERSION)

down:
	docker compose down
	
clean:
	docker compose down -v

up:
	docker compose up --build -d