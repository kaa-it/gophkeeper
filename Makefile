.PHONY: build

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