VERSION := $(shell cat Version)
BINARY_NAME := netron
BUILD_DIR := build

.PHONY: all clean build-all build-linux-amd64 build-linux-arm64 build-linux-armv7 build-linux-armv6 build-linux-i386

all: build-all

clean:
	rm -rf $(BUILD_DIR)

build-all: clean
	@echo "Building version $(VERSION)"
	mkdir -p $(BUILD_DIR)
	$(MAKE) build-linux-amd64
	$(MAKE) build-linux-arm64
	$(MAKE) build-linux-armv7
	$(MAKE) build-linux-armv6
	$(MAKE) build-linux-i386

build-linux-amd64:
	@echo "Building for Linux AMD64..."
	mkdir -p $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64
	cp -r static $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64/
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64/$(BINARY_NAME) .
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-$(VERSION)-linux-amd64

build-linux-arm64:
	@echo "Building for Linux ARM64..."
	mkdir -p $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64
	cp -r static $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64/
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64/$(BINARY_NAME) .
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-$(VERSION)-linux-arm64

build-linux-armv7:
	@echo "Building for Linux ARMv7..."
	mkdir -p $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-armv7
	cp -r static $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-armv7/
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-armv7/$(BINARY_NAME) .
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-armv7.tar.gz $(BINARY_NAME)-$(VERSION)-linux-armv7

build-linux-armv6:
	@echo "Building for Linux ARMv6..."
	mkdir -p $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-armv6
	cp -r static $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-armv6/
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-armv6/$(BINARY_NAME) .
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-armv6.tar.gz $(BINARY_NAME)-$(VERSION)-linux-armv6

build-linux-i386:
	@echo "Building for Linux i386..."
	mkdir -p $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-i386
	cp -r static $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-i386/
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-i386/$(BINARY_NAME) .
	cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)-$(VERSION)-linux-i386.tar.gz $(BINARY_NAME)-$(VERSION)-linux-i386

dev:
	go run . --run -p 8080

install-deps:
	go mod tidy

test:
	go test ./...