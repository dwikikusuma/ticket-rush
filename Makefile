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

.PHONY: proto

PROTO_DIR := common/proto
MODULE := github.com/dwikikusuma/ticket-rush

proto-tools:
	@command -v protoc >/dev/null 2>&1 || (echo "protoc is not installed" && exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || (echo "missing protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" && exit 1)
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || (echo "missing protoc-gen-go-grpc: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" && exit 1)

proto: proto-tools
	protoc -I $(PROTO_DIR) \
		--go_out=. --go_opt=module=$(MODULE) --go_opt=paths=import \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) --go-grpc_opt=paths=import \
		$$(find $(PROTO_DIR) -name "*.proto")
