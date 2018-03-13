# scccmd - Spring Cloud Config cli tool

[![Build Status](https://travis-ci.org/WanderaOrg/scccmd.svg?branch=master)](https://travis-ci.org/WanderaOrg/scccmd)
[![Docker Build Status](https://img.shields.io/docker/build/wanderadock/scccmd.svg)](https://hub.docker.com/r/fluffycroc/scccmd/)
[![Go Report Card](https://goreportcard.com/badge/github.com/WanderaOrg/scccmd)](https://goreportcard.com/report/github.com/WanderaOrg/scccmd)
[![GitHub release](https://img.shields.io/github/release/WanderaOrg/scccmd.svg)](https://github.com/WanderaOrg/scccmd/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/WanderaOrg/scccmd/blob/master/LICENSE)

Tool for obtaining configuration from config server

### How to develop
* Checkout into your GOROOT directory (e.g. /go/src/github.com/WanderaOrg/scccmd)
* `cd` into the folder and run `dep ensure`
* Tests are started by `go test -v ./...`
* Or if you dont want to setup your local go env just use the provided Dockerfile

### Docker repository
Tool is released as docker image as well, check the [repository](https://hub.docker.com/r/wanderadock/scccmd/).

### Tool documentation
[docs](docs/config.md)	 - Generated documentation for the tool