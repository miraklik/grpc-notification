.PHONY run migrations migrate-up migrate-down

migrations:
	migrate create -ext sql -dir internal/db/migrations -seq create_notification_table

migrate-up:
	migrate -path internal/db/migrations -database <DB_URL> up

migrate-down:
	migrate -path internal/db/migrations -database <DB_URL> down

run:
	go run cmd/main.go