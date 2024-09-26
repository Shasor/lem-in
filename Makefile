all: lem_in visualizer

lem_in:
	@echo "Build lem-in in progress..."
	@go build -o lem-in cmd/lem-in/main.go
	@echo "OK!"

visualizer:
	@echo "Build visualizer in progress..."
	@go build -o visualizer cmd/visualizer/main.go
	@echo "OK!"