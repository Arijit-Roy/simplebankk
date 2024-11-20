createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

migrateup:
	migrate -path db/migrations/ -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migrations/ -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose down

sqlc:
	./sqlc generate

test:
	go test -v -cover ./...
server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

# docker build -t simplebank:latest .
# docker run --name simplebank --network bank-network -p 8080:8080 \ 
# -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" \ 
# -e SERVER_ADDRESS=0.0.0.0:8080 simplebank:latest

.PHONY: createdb  dropdb postgres migrateup migratedown sqlc test server mock