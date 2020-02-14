package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wandera/scccmd/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
