migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:admin@localhost:5432/ecom?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:admin@localhost:5432/ecom?sslmode=disable" -verbose down

.PHONY: migrate-up migrate-down
