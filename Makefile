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