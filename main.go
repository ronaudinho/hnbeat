package main

import (
	"os"

	"github.com/ronaudinho/hnbeat/cmd"

	_ "github.com/ronaudinho/hnbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
