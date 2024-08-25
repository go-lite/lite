build:
	go build -v ./...

analyze:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

test:
	go test ./...

cover:
	@GOEXPERIMENT=nocoverageredesign go test -race -coverprofile=coverage.out -covermode=atomic ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1 run --config scripts/.golangci.yaml

lint-fix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1 run --config scripts/.golangci.yaml --fix

.PHONY: build test cover lint lint-fix fix
