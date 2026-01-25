INSTALL_DIR := $(HOME)/.local/bin
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(HOME)/go/bin
endif

.PHONY: install uninstall

install:
	@echo "Installing Go tools to $(GOBIN)..."
	@find . -mindepth 2 -name go.mod -exec dirname {} \; | while read dir; do \
		echo "  Installing $$(basename $$dir)..."; \
		(cd $$dir && go install .); \
	done
	@mkdir -p $(INSTALL_DIR)
	@echo "Installing scripts to $(INSTALL_DIR)..."
	@find bin -maxdepth 1 -type f -perm +111 | while read f; do \
		name=$$(basename $$f); \
		ln -sf $(CURDIR)/$$f $(INSTALL_DIR)/$$name; \
		echo "  Linked $$name"; \
	done

uninstall:
	@echo "Uninstalling Go tools from $(GOBIN)..."
	@find . -mindepth 2 -name go.mod -exec dirname {} \; | while read dir; do \
		tool_name=$$(basename $$dir); \
		rm -f $(GOBIN)/$$tool_name; \
		echo "  Removed $$tool_name"; \
	done
	@echo "Uninstalling scripts from $(INSTALL_DIR)..."
	@find bin -maxdepth 1 -type f -perm +111 | while read f; do \
		name=$$(basename $$f); \
		rm -f $(INSTALL_DIR)/$$name; \
		echo "  Removed $$name"; \
	done
