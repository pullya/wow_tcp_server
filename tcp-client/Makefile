.PHONY: run_client
run_client:
	go run ./cmd/tcp-client/main.go

.PHONY: run
run:
	go run ./cmd/tcp-client/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: .test
.test:
	$(info Running tests...)
	go test ./...

.PHONY: test
test: .test