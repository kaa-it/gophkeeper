.PHONY: build race msan pre-commit run-pre-commit lint gofumpt goimports tools

CLIENT_VERSION := 0.1.0
SERVER_VERSION := 0.1.0

COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +'%Y/%m/%d %H:%M:%S')

build:
	go build -o gophkeeper_client -ldflags \
		"-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildVersion=${CLIENT_VERSION} \
		-X 'github.com/kaa-it/gophkeeper/pkg/buildconfig.buildDate=${DATE}' \
		-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildCommit=${COMMIT}" ./cmd/client ;
	go build -o gophkeeper_server -ldflags \
		"-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildVersion=${SERVER_VERSION} \
		-X 'github.com/kaa-it/gophkeeper/pkg/buildconfig.buildDate=${DATE}' \
		-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildCommit=${COMMIT}" ./cmd/server ;

pre-commit:
	pre-commit install

run-pre-commit:
	pre-commit run --all-files

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@v0.6.0

lint:
	golangci-lint run ./internal/...

gofumpt:
	gofumpt -l -w ./internal

goimports:
	goimports -w --local github.com/kaa-it/gophkeeper ./internal