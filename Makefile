BINARY_AMD := kiwipanel-amd
BINARY_ARM := kiwipanel-arm
ZIP_NAME := kiwipanel.zip
SRC_DIR := kiwipanel

.PHONY: build-amd build-arm zip
run: build
	@sudo ./bin/kiwipanel start
build:
	go build -o kiwipanel ./cmd/kiwipanel
build-amd:	
	GOOS=linux GOARCH=amd64 go build -o kiwipanel/bin/kiwipanel ./cmd/kiwipanel
build-arm:
	GOOS=linux GOARCH=arm64 go build -o kiwipanel/bin/kiwipanel ./cmd/kiwipanel
# Zip the folder including any binaries present
zip:
	@echo "Removing existing $(ZIP_NAME) if it exists..."
	@rm -f $(ZIP_NAME)
	@echo "Creating $(ZIP_NAME)..."
	# Include kiwipanel/ folder plus binaries if they exist
	zip -r $(ZIP_NAME) $(SRC_DIR) $(wildcard $(BINARY_AMD) $(BINARY_ARM)) -x "*.DS_Store" -x "__MACOSX/*"
	@echo "$(ZIP_NAME) created successfully."
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