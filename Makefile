include .env

.PHONY: build clean

APP_NAME := build/backup_app

build:
	go build -ldflags="-s -w" -o $(APP_NAME) cmd/backup_app.go

clean:
	rm -f $(APP_NAME)

run:
	./build/backup_app --backup --restore --docker \
    --ssh_host=$(SSH_HOST) \
    --ssh_port=$(SSH_PORT) \
    --ssh_user=$(SSH_USER) \
    --ssh_password=$(SSH_PASSWORD) \
    --db_name=$(DB_NAME) \
    --db_user=$(DB_USER) \
    --db_password=$(DB_PASSWORD) \
    --db_host=$(DB_HOST) \
	--restore_db_name=$(RESTORE_DB_NAME) \
	--restore_db_user=$(RESTORE_DB_USER) \
	--restore_db_password=$(RESTORE_DB_PASSWORD) \
	--restore_db_host=$(RESTORE_DB_HOST)