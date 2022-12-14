package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
- How to create a basic logging middleware in Go.
- A middleware simply takes a http.HanderFunc as one of its parameters, wraps it and returns a new
  http.HandlerFunc for the server to call.
*/

func logging(f http.HandlerFunc) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    log.Println(r.URL.Path)
    f(w, r)
  }
}

func foo(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "foo")
}

func bar(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "bar")
}

func main() {
  http.HandleFunc("/foo", logging(foo))
  http.HandleFunc("/bar", logging(bar))

  http.ListenAndServe(":8081", nil)
}
