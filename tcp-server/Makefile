GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

.PHONY: run_server
run_server:
	go run ./cmd/tcp-server/main.go

.PHONY: run
run:
	go run ./cmd/tcp-server/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: generate
generate:
	$(BUF_EXE) generate

# .PHONY: .build
# build: 
# 	generate build

.PHONY: build
.build:
	CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-o ./bin/tcp-server$(shell go env GOEXE) ./cmd/tcp-server/main.go