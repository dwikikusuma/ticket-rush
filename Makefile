# Database Connection String
DB_URL=postgres://user:password@localhost:5432/ticket_db?sslmode=disable

# 1. Migration Commands
.PHONY: migrate-up migrate-down migrate-force

# Run migrations (UP)
migrate-up:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" up

# Rollback migrations (DOWN)
migrate-down:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" down 1

# Force version (useful if migration gets dirty)
migrate-force:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" force 1

# 2. Seeder Command
.PHONY: seed
seed:
	go run cmd/seeder/main.go
