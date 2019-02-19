package main

import (
	"github.com/WanderaOrg/scccmd/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
