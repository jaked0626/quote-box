include app.env

serve:
	go run ./cmd/web -addr=${HTTP_TARGET_ADDRESS} -dsn=${DB_SOURCE}

serve_log:
	go run ./cmd/web -addr=${HTTP_TARGET_ADDRESS} -dsn=${DB_SOURCE} >>./tmp/info.log 2>>./tmp/error.log

container:
	docker run --name snippet_pg -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:14-alpine

startpg:
	docker exec -it snippet_pg psql -U ${DB_USER} ${DB_NAME}

createdb:
	docker exec -it snippet_pg createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

# grants permission to all tables in database. Wait until after all tables are created with migration. 
createdbrole:
	- docker exec -it snippet_pg psql -U ${DB_USER} -c "CREATE USER ${DB_ROLE} WITH PASSWORD '${DB_ROLE_PW}';"
	- docker exec -it snippet_pg psql -U ${DB_USER} -c "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO ${DB_ROLE};"

migrateup: 
	- migrate -path ./internal/migration/ -database ${DB_SOURCE} -verbose up

migratedown: 
	- migrate -path ./internal/migration/ -database ${DB_SOURCE} -verbose down

up: 
	- make container 
	- sleep 2
	- make createdb
	- sleep 2
	- make createdbrole
	- sleep 2
	- make migrateup
	- sleep 2
	- make serve

down:
	- make migratedown
	- sleep 2
	- docker stop snippet_pg && docker rm -v snippet_pg
	
dropdb:
	docker exec -it snippet_pg dropdb ${DB_NAME}

.PHONY: serve serve_log startpg container createdb createdbrole dropdb up down migrateup migratedown