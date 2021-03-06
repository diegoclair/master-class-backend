DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 --network simple-bank-local -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

createmigration:
#how to call example: filename=t make createmigration
#withou @ before if, it logs the entire make command
	@if [ "$(filename)" = "" ]; then\
        echo "filename parameter is required";\
	else \
		migrate create -ext sql -dir db/migrations -seq "$(filename)";\
    fi

dropdb: 
	docker exec -it postgres12 dropdb simple_bank

migratedown: 
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migratedownlast: 
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

migrateup: 
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrateuplast: 
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1


mock:
# -package = name of package in the generated file;
# -destination = where the file will be generated and with name
# Store = interface that we want to generate the mock, we can pass more than one with comma
	mockgen -package mockdb -destination db/mock/store.go github.com/diegoclair/master-class-backend/db/sqlc Store

server:
	go run main.go

# after we execute sqlc generate it will update the querier interface and so we need to regenerate the mock again
sqlc: sqlccommand mock

sqlccommand:
#this command generate code for the queries defined inside of db/query
	sqlc generate 

test:
	go test -v -cover ./...

proto:
	rm -f pb/*.go
	rm -f docs/*.swagger.json
	protoc --proto_path=proto/protodefs --go_out=proto/pb --go_opt=paths=source_relative \
    --go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto/pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=docs --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/protodefs/*.proto

evans:
	evans -p 9000 --host localhost -r repl

.PHONY: postgres createdb createmigration dropdb migratedown migratedownlast migrateup migrateuplast sqlc test server mock proto evans