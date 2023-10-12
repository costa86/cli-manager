.DEFAULT_GOAL := run
BINARY=cli-manager
BUILD_CMD := go build --trimpath -ldflags="-s -w"
BUILD_CMD_LINUX := go build -ldflags="-s -w"

run:
	@go run *.go

build:
	@go mod tidy

	@echo "Building for Linux"
	@GOOS=linux GOARCH=amd64 $(BUILD_CMD_LINUX) -o $(BINARY)-linux.exe

	@echo "Building for Windows" 
	@GOOS=windows GOARCH=amd64 $(BUILD_CMD) -o $(BINARY)-windows.exe

	@echo "Building for MacOS"
	@GOOS=darwin GOARCH=amd64 $(BUILD_CMD) -o $(BINARY)-darwin