DB_URL=postgresql://root:secret@localhost:5432/account?sslmode=disable

network:
	docker network create foundation

postgres:
	docker run --name account-postgres --network foundation -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it account-postgres createdb --username=root --owner=root account

dropdb:
	docker exec -it account-postgres dropdb account --username=root 

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate%:
	migrate -path db/migration -database "$(DB_URL)" -verbose up $*

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build docs/template.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/template.dbml
	cp docs/schema.sql db/migration/000000_init_schema.up.sql

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

proto:
	rm -f docs/swagger/*.swagger.json
	rm -rf gen/pb/*
	buf generate
	statik -src=./docs/swagger -dest=./docs

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis