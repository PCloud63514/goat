PROJECT_NAME := goat
PKG_PATH_SERVER := ./cmd

verify:
	go mod tidy
	go test -v ./...

build:
	go mod tidy
	go build -o $(PROJECT_NAME) $(PKG_PATH_SERVER)
