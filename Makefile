run: build
	@sudo ./bin/kiwipanel start
build:
	go build -o kiwipanel ./cmd/kiwipanel
build-zip: build
	zip -j kiwipanel.zip ./bin/kiwipanel ./configs/sample_config.yaml
build-amd:	
	GOOS=linux GOARCH=amd64 go build -o bin/kiwipanel ./cmd/kiwipanel
build-arm:
	GOOS=linux GOARCH=arm64 go build -o bin/kiwipanel ./cmd/kiwipanel
dev:
	go run ./cmd/kiwipanel dev
start:
	@sudo systemctl start kiwipanel
status:
	@sudo systemctl status kiwipanel
restart:
	@sudo systemctl restart kiwipanel
stop:
	@sudo systemctl stop kiwipanel
## css: build tailwindcss we use tailwindcss cli ./tailwindcss
css:
	./tailwindcss -i ./assets/css/tailwind.css -o ./assets/css/main.css --minify
## css-watch: watch build tailwindcss 
watch:
	./tailwindcss -i ./assets/css/tailwind.css -o ./assets/css/main.css --watch