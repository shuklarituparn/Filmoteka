setup:
	@echo "=== Installing Tools ==="
	@bash ./scripts/install_tools.sh
	@echo "=== Tools Installed ==="

lint:
	@echo "=== Running Linter ==="
	@./bin/golangci-lint run
	@echo "=== Linter Completed ==="

