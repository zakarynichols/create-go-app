# Create Go App

Create a simple, minimal, and unopinionated go app with a single command.

To create an application called 'my-app', use `go run` to compile and execute the program remotely.

```sh
$ go run github.com/zakarynichols/create-go-app my-app
$ cd ./my-app
$ go run main.go
> Listening on port 9999
```

In a new terminal, ping the running web server.

```sh
$ curl localhost:9999/foo
> path: /foo
$ curl localhost:9999/bar
> path: /bar
```
