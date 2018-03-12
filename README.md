# scccmd - Spring Cloud Config cli tool

[![Build Status](https://travis-ci.org/WanderaOrg/scccmd.svg?branch=master)](https://travis-ci.org/WanderaOrg/scccmd)

Tool for obtaining configuration from config server

### How to develop
* Checkout into your GOROOT directory (e.g. /go/src/github.com/wanderaorg/scccmd)
* `cd` into the folder and run `dep ensure`
* Tests are started by `go test -v ./...`
* Or if you dont want to setup your local go env just use the provided Dockerfile

### Tool documentation
[docs](docs/config.md)	 - Generated documentation for the tool