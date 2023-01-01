package code

// HTTP is the template for main.go when initializing an http server.
const HTTP string = `
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const Port = 9999

func main() {
	foo := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("path: %s\n", r.URL.Path))
	}

	bar := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("path: %s\n", r.URL.Path))
	}

	fmt.Printf("Listening on port %d\n", Port)

	http.HandleFunc("/foo", foo)
	http.HandleFunc("/bar", bar)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}
`
