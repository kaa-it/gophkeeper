.PHONY: build race msan pre-commit run-pre-commit lint gofumpt goimports tools proto

CLIENT_VERSION := 0.1.0
SERVER_VERSION := 0.1.0

COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +'%Y/%m/%d %H:%M:%S')
PKG_LIST := $(shell go list ./internal/... ;)

build_client:
	go build -o gophkeeper_client -ldflags \
		"-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildVersion=${CLIENT_VERSION} \
		-X 'github.com/kaa-it/gophkeeper/pkg/buildconfig.buildDate=${DATE}' \
		-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildCommit=${COMMIT}" ./cmd/client ;

build_server:
	go build -o gophkeeper_server -ldflags \
		"-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildVersion=${SERVER_VERSION} \
		-X 'github.com/kaa-it/gophkeeper/pkg/buildconfig.buildDate=${DATE}' \
		-X github.com/kaa-it/gophkeeper/pkg/buildconfig.buildCommit=${COMMIT}" ./cmd/server ;

proto:
	mkdir -p ./internal/pb
	protoc --proto_path=internal/proto --go_out=./internal/pb --go_opt=paths=source_relative \
	  --go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
	  internal/proto/*.proto

pre-commit:
	pre-commit install

run-pre-commit:
	pre-commit run --all-files

tools:
	#curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@v0.6.0
	go install github.com/segmentio/golines@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/ktr0731/evans@latest

lint:
	golangci-lint run ./internal/...

gofumpt:
	gofumpt -l -w ./internal

goimports:
	goimports -w --local github.com/kaa-it/gophkeeper ./internal

evans:
	evans -r repl -p 8080

cert:
	cd cert; ./gen.sh; cd ..

run_server:
	./gophkeeper_server -port 8080

run_client:
	./gophkeeper_client -address "0.0.0.0:8080"

test:
	go test ./internal/...