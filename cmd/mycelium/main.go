package main

import (
	"os"

	"github.com/robertguss/mycelium/internal/cli"
)

func main() {
	os.Exit(cli.Main(os.Args))
}
