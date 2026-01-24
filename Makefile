BIN_DIR := $(CURDIR)/bin
INSTALL_DIR := $(HOME)/.local/bin

.PHONY: all build scripts clean install uninstall

all: build scripts

# Builds everything and links to local ./bin folder
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

# Links everything in the ./bin folder to the BIN_DIR (~/.local/bin)
install: all
	@mkdir -p $(INSTALL_DIR)
	@echo "Installing to $(INSTALL_DIR)..."
	@for f in $(BIN_DIR)/*; do \
		name=$$(basename $$f); \
		ln -sf $$f $(INSTALL_DIR)/$$name; \
		echo "  Linked $$name"; \
	done

uninstall:
	@echo "Uninstalling from $(INSTALL_DIR)..."
	@for f in $(BIN_DIR)/*; do \
		name=$$(basename $$f); \
		rm -f $(INSTALL_DIR)/$$name; \
		echo "  Removed $$name"; \
	done

clean:
	rm -rf $(BIN_DIR)
