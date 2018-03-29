package main

import (
	"github.com/WanderaOrg/scccmd/cmd"
	"os"
)

func main() {

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}

}
