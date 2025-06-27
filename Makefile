.PHONY: lint
lint:
	@echo "==> Lint <=="
	if [ ! -e ./bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.1.6; fi
	./bin/golangci-lint run ./... --timeout=10m

.PHONY: format
format:
	@echo "==> Format <=="
	@echo "--> Formatting code"
	find . -name '*.go' -exec gofumpt -w {} +;
	@echo "--> Formatting imports"
	find . -name '*.go' -exec goimports -w -local=github.com/mats852/chip {} +;
