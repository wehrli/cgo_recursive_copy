BINARY_NAME=copy_recursive_bin

GO=go

GOFLAGS=

SOURCES=$(wildcard *.go)

OUTPUT_DIR=bin

OUTPUT_BINARY=$(OUTPUT_DIR)/$(BINARY_NAME)

BUILD_FLAGS=

CLEAN_FLAGS=-r

.DEFAULT_GOAL: all

$(OUTPUT_BINARY): $(SOURCES)
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) $(BUILD_FLAGS) -o $(OUTPUT_BINARY) $(SOURCES)

clean:
	@echo "Cleaning..."
	@rm $(CLEAN_FLAGS) $(OUTPUT_DIR)

all: clean $(OUTPUT_BINARY)
