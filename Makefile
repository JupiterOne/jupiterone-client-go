TEST?=$$(go list ./...)
PKG_NAME=jupiterone

default: build

build: fmtcheck
	go install

test: fmtcheck
	JUPITERONE_API_KEY=fake JUPITERONE_ACCOUNT_ID=fake RECORD=false TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout=5m -parallel=4

fmpt: 
	gofumpt -l -w .

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./$(PKG_NAME)

fmtcheck:
	@./scripts/fmtcheck.sh

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(PKG_NAME)/...

ok:
	@echo OK

precommit: tidy fmt lint test sec ok

tidy: 
	go mod tidy

tools:
	@echo "==> installing required tooling..."
	go install github.com/client9/misspell/cmd/misspell
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc cassettes fmtcheck lint tools test-compile docs
