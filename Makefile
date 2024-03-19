.PHONY: help
help:
	@echo "Available commands:"
	@awk '/^[a-zA-Z0-9_-]+:.*?##/ {split($$0, a, ": ## "); printf "  %-20s %s\n", a[1], a[2]}' $(MAKEFILE_LIST)

## example: make cm name=test
create-migration: ## Create migration file in migration directory using goose, use name variable as an input for sql migration name 
	@echo "Creating migration file..."
	@goose -dir db/migrations create ${name} sql

## example: make cs name=test
create-seeder: ## Create seeder file in seeder directory using goose, use name variable as an input for sql seeder name 
	@echo "Creating migration file..."
	@goose -dir db/seeders create ${name} sql

up-migration: ## Migrate up migration file into the database using goose
	@echo 'Migrating file into database...'
	@goose -dir ./db/migrations postgres "postgresql://skeleton:skeleton@localhost:5432/skeletonDB?sslmode=disable" up

up-seeder: ## Seed the database using goose
	@echo 'Seeding file into database...'
	@goose -dir ./db/seeders -no-versioning postgres "postgresql://skeleton:skeleton@localhost:5432/skeletonDB?sslmode=disable" up

## docker exec -it <container_name> psql -U <username> -d <database>
exec-db: ## Exec into the postgres container
	@echo 'Executing into postgres container...'
	@docker exec -it skeleton_db psql -U skeleton -d skeletonDB

## docker exec -it <container_name> /bin/bash
exec-kafka: ## Exec into the kafka container
	@echo 'Executing into kafka container...'
	@docker exec -it skeleton_kafka /bin/bash

compose-up: ## Run docker-compose
	@echo 'Running docker-compose...'
	@docker-compose up --build

compose-down: ## Run docker-compose down
	@echo 'Running docker-compose down...'
	@docker-compose down -v

sqlc: ## Generate query into golang using sqlc
	@echo 'Generating query...'
	@sqlc generate

run: up-migration ## Run the server
	@echo 'Running the server...'
	@go run cmd/server/main.go

k6: ## Run k6 to test the server
	@k6 run script/test.js --out influxdb=http://localhost:8086/k6