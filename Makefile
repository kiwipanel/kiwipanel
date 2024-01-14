run: build
	@ ./bin/kiwipanel start
build:
	@go build -o bin/kiwipanel