package main

import (
	"github.com/wanderaorg/scccmd/config/cmd"
	"os"
)

func main() {

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}

}
