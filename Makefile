build:
	go build -v ./...

analyze:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

test:
	go test ./...

cover:
	@GOEXPERIMENT=nocoverageredesigngo test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run --config scripts/.golangci.yaml

lint-fix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run --config scripts/.golangci.yaml --fix

.PHONY: build test cover lint lint-fix fix