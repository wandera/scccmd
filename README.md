# scccmd - Spring Cloud Config cli tool

[![Test](https://github.com/wandera/scccmd/actions/workflows/test.yml/badge.svg)](https://github.com/wandera/scccmd/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wandera/scccmd)](https://goreportcard.com/report/github.com/wandera/scccmd)
[![GitHub release](https://img.shields.io/github/release/wandera/scccmd.svg)](https://github.com/wandera/scccmd/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/wandera/scccmd/blob/master/LICENSE)

Tool for obtaining configuration from config server

### How to develop
* Run `go get github.com/wandera/scccmd`
* Build by `go build -v`
* Tests are started by `go test -v ./...`
* Or if you dont want to setup your local go env just use the provided Dockerfile

### Docker repository
The tool is released as docker image as well, check the [repository](https://github.com/wandera/git2kube/pkgs/container/scccmd).

### Kubernetes Initializer
The tool could be used as Webhook for Kubernetes deployments. 
Deployed webhook will add init container to applicable deployments,
which in turn downloads configuration in deployment initialization phase.
Example k8s [manifest](docs/k8s/bundle.yaml).

### Tool documentation
[docs](docs/scccmd.md)	 - Generated documentation for the tool
