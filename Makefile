include app.env

serve:
	go run ./cmd/web -addr=${HTTP_TARGET_ADDRESS}

postgres:
	docker run --name postgres -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${DB_ROOT_USER} -e POSTGRES_PASSWORD=${DB_ROOT_PW} -d postgres:14-alpine

createdb:
	docker exec -it postgres createdb --username=${DB_ROOT_USER} --owner=${DB_ROOT_USER} ${DB_NAME}

startpg:
	docker exec -it postgres psql -U ${DB_ROOT_USER} ${DB_NAME}

# grants permission to all tables in database. Wait until after all tables are created with migration. 
createdbrole:
	docker exec -it postgres psql -U ${DB_ROOT_USER} -c "CREATE USER ${DB_USER} WITH PASSWORD '${DB_USER_PW}';"

grantpermission:
	docker exec -it postgres psql -U ${DB_ROOT_USER} -c "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO ${DB_USER};"

dropdb:
	docker exec -it postgres dropdb ${DB_NAME}

rmpostgres:
	docker stop postgres && docker rm -v postgres

init_migrate:
	migrate create -ext sql -dir ./internal/db/migration -seq init_schema

migrateup:
	migrate -path internal/db/migration -database "postgresql://root:secret@localhost:5432/snippet_db?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://root:secret@localhost:5432/snippet_db?sslmode=disable" -verbose down

sqlc: 
	sqlc generate ./internal/db/erd

.PHONY: serve startpg postgres createdb dropdb rmpostgres migrateup migratedown