BIN_DIR := $(CURDIR)/bin

.PHONY: all build scripts clean

all: build scripts

build:
	@mkdir -p $(BIN_DIR)
	@echo "Building Go tools..."
	@find . -mindepth 2 -name go.mod -exec dirname {} \; | while read dir; do \
		tool_name=$$(basename $$dir); \
		echo "  $$tool_name"; \
		(cd $$dir && go build -o $(BIN_DIR)/$$tool_name .); \
	done

scripts:
	@mkdir -p $(BIN_DIR)
	@echo "Linking scripts..."
	@find . -maxdepth 1 -type f -perm +111 -not -name "Makefile" -not -path "./.git/*" -exec ln -sf $(CURDIR)/{} $(BIN_DIR)/ \;

clean:
	rm -rf $(BIN_DIR)
