include .env

serve:
	go run ${MAIN_DIRECTORY} -addr=${HTTP_TARGET_ADDRESS}

postgres:
	docker run --name postgres -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:14-alpine

startpg:
	docker exec -it postgres psql -U ${DB_USER} ${DB_NAME}

createdb:
	docker exec -it postgres createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

dropdb:
	docker exec -it postgres dropdb ${DB_NAME}

rmpostgres:
	docker stop postgres && docker rm -v postgres

.PHONY: serve startpg postgres createdb dropdb rmpostgres