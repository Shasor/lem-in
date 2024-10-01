all: lem_in visualize

lem_in:
	@echo "Build lem-in in progress..."
	@go build -o lem-in cmd/lem-in/main.go
	@echo "OK!"

visualize:
	@echo "Build visualizer in progress..."
	@go build -o visualizer cmd/visualizer/main.go
	@echo "OK!"