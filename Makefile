text =

migratecreate:
	migrate create -ext sql -dir internal/adapter/storage/postgres/migrations -seq $(text) 