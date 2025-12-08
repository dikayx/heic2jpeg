# Go settings
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
BINARY_NAME := heic2jpeg
TARGET_DIR := target
CMD_PATH := ./cmd/heic2jpeg

# Ensure target directory exists
$(TARGET_DIR):
	mkdir -p $(TARGET_DIR)

# Default target
.PHONY: all
all: build

# Build the binary into target/
.PHONY: build
build: $(TARGET_DIR)
	$(GOBUILD) -o $(TARGET_DIR)/$(BINARY_NAME) $(CMD_PATH)

# Run directly with Go (development mode)
.PHONY: run
run:
	$(GORUN) $(CMD_PATH)

# Run tests
.PHONY: test
test:
	$(GOTEST) ./...

# Clean up binaries in target/
.PHONY: clean
clean:
	rm -rf $(TARGET_DIR)

# Cross compile for multiple OS/ARCH into target/
.PHONY: release
release: $(TARGET_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(TARGET_DIR)/$(BINARY_NAME).exe $(CMD_PATH)
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o $(TARGET_DIR)/$(BINARY_NAME)_mac $(CMD_PATH)
	GOOS=linux   GOARCH=amd64 $(GOBUILD) -o $(TARGET_DIR)/$(BINARY_NAME)_linux $(CMD_PATH)
