.PHONY: migrate-up migrate-down migrate-create

# Database migration commands
migrate-up:
	migrate -database "${DATABASE_URL}" -path db/migrations up

migrate-down:
	migrate -database "${DATABASE_URL}" -path db/migrations down

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name