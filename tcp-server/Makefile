.PHONY: run_server
run_server:
	go run ./cmd/tcp-server/main.go

.PHONY: run
run:
	go run ./cmd/tcp-server/main.go

.PHONY: lint
lint:
	golangci-lint run ./...
