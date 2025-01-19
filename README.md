# Create Go App

_This project is still under development. Many features are not available and subject to change._

[![Go Reference](https://pkg.go.dev/badge/create-go-app.com.svg)](https://pkg.go.dev/create-go-app.com)

Create a production-ready, full-stack application, with sensible defaults.

## HTTP Server

`$ go run create-go-app.com@latest -http my-http-server`

## CLI

_Creating a command line interface in one command is still under development._

`$ go run create-go-app.com@latest -cli my-cli`

## Development

Make sure to `god mod init` and `go get` in `create-go-app/emit`. This will prevent compile time errors. Auto-generated `go.sum` and `go.mod` are ignored by source control.