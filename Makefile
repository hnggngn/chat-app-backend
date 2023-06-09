.PHONY: build run

build:
	swag fmt && swag init -g cmd/app/main.go --ot "go"

run: build
	air