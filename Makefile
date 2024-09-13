all: build

build:
	@echo "Build in progress..."
	@go build -o lem-in cmd/lem-in/main.go
	@echo "OK!"