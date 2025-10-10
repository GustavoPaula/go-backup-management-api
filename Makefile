text =

migratecreate:
	migrate create -ext sql -dir internal/adapter/storage/postgres/migrations -seq $(text) 

migraterollback:
	migrate -path ./internal/adapter/storage/postgres/migrations -database "postgres://postgres:docker@localhost:5432/db?sslmode=disable" down 1
