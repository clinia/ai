GO := go

EXAMPLE_MAIN_FILES := $(shell find examples -name main.go)
EXAMPLE_DIRS := $(sort $(patsubst %/,%,$(dir $(EXAMPLE_MAIN_FILES))))

.PHONY: test-examples
test-examples:
	@set -e; \
	failed_examples=(); \
	for dir in $(EXAMPLE_DIRS); do \
		echo "==> $(GO) run ./$$dir"; \
		if ! $(GO) run ./$$dir; then \
			failed_examples+=("$$dir"); \
		fi; \
	done; \
	if [ $${#failed_examples[@]} -eq 0 ]; then \
		echo "All examples passed."; \
	else \
		echo ""; \
		echo "Example tests failed. The following examples failed:"; \
		for failed in "$${failed_examples[@]}"; do \
			echo "  - $$failed"; \
		done; \
		exit 1; \
	fi
