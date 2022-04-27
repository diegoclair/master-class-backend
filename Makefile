postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it postgres12 dropdb simple_bank

migrateup: 
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
# -package = name of package in the generated file;
# -destination = where the file will be generated and with name
# Store = interface that we want to generate the mock, we can pass more than one with comma
	mockgen -package mockdb -destination db/mock/store.go github.com/diegoclair/master-class-backend/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc, test, server, mock