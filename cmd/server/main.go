package main

import (
	"github.com/kaa-it/gophkeeper/internal/server"
	"github.com/kaa-it/gophkeeper/pkg/buildconfig"
)

func main() {
	buildconfig.PrintBuildInfo()

	server.Run()
}
