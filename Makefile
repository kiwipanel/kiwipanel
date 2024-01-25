run: build
	@sudo ./bin/kiwipanel start
build:
	@go build -o bin/kiwipanel
dev: build
	@ ./bin/kiwipanel dev
start:
	@sudo systemctl start kiwipanel
status:
	@sudo systemctl status kiwipanel
stop:
	@sudo systemctl stop kiwipanel