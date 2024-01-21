run: build
	@ ./bin/kiwipanel start
build:
	@go build -o bin/kiwipanel
dev: build
	@ ./bin/kiwipanel dev