run: build
	@sudo ./bin/kiwipanel start
#ARM: install: line 368: /usr/local/go/bin/go: cannot execute binary file: Exec format error
build:
	@go build -o bin/kiwipanel
dev: build
	@ ./bin/kiwipanel dev
start:
	@sudo systemctl start kiwipanel
status:
	@sudo systemctl status kiwipanel
restart:
	@sudo systemctl restart kiwipanel
stop:
	@sudo systemctl stop kiwipanel