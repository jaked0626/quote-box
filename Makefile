include .env

serve:
	go run ${MAIN_DIRECTORY} -addr=${HTTP_TARGET_ADDRESS}

.PHONY: serve