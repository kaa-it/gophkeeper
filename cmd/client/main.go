package main

import (
    "github.com/kaa-it/gophkeeper/internal/client"
    "github.com/kaa-it/gophkeeper/pkg/buildconfig"
)

func main() {
	buildconfig.PrintBuildInfo()

	client.Run()
}