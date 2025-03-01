.PHONY: build clean

APP_NAME := build/backup_app

build:
	go build -ldflags="-s -w" -o $(APP_NAME) cmd/backup_app.go

clean:
	rm -f $(APP_NAME)