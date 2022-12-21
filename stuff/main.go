package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Create Go App")
	http.ListenAndServe(":1337", nil)
}
