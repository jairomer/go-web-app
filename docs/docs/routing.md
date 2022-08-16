# Routing

One thing that `net/http` does not do very well is complex request routing like segmenting a request url into single parameters.

`gorilla/mux` is a package that solves this issue and provides a lot of features to increase productivity when writing web applications.
It is also compliant with Go's default request handler signature.

It can be installed with `go get -u github.com/gorilla/mux`.

## Creating a new Router

- The router is the main router for your web application and will later be passed as parameter to the server.
- It will receive all HTTP connections and pass it on to the request handlers you will register on it.

```go
r := mux.NewRouter()
```

## Registering a Request Handler

- Once you have a new router you can register request handlers like usual.
- The only difference is that instead of calling `http.HandleFunc(...)` you will use `r.HandleFunc(...)`.

## URL Parameters

- The biggest strength of the `gorilla/mux` router is the ability to extract segments from the request URL.
- Example: `/books/go-programming-blueprint/page/10`
- This URL has two dynamic segments.
  + Book title slig
  + Page number
- To match the request handler match the URL mentioned above, you will need to use placeholders in the pattern.
- To get the data on these segments, use the function `mux.Vars(r)`.
```go
r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
  // get the book variables from the request.
  vars := mux.Vars(r)
  title := vars["title"]
  page := vars["page"]
  // navigate to the page.
})
```

## Setting the HTTP server's router

- The second parameter for `http.ListenAndServe(":80", nil)` is the router of the HTTP server.
- `nil` references the default router.
- `http.ListenAndServe(":80", r)`

## Features

### Methods

```go
// Restrict the handler to specific methods.
r.HandleFunc("/books/{title}", CreateBook).Methods("POST")
r.HandleFunc("/books/{title}", ReadBook).Methods("GET")
r.HandleFunc("/books/{title}", UpdateBook).Methods("PUT")
r.HandleFunc("/books/{title}", DeleteBook).Methods("DELETE")
```

### Hostnames and Subdomains

```go
// Restrict the handler to specific hostnames or subdomains.
r.HandleFunc("/books/{title}", BookHandler).Host("www.mybookstore.com")
```

## Schemes

```go
// Restrict the request handler to http/https.
r.HandleFunc("/secure", SecureHandler).Schemes("https")
r.HandleFunc("/insecure", InsecureHandler).Schemes("http")
```

## Path Prefixes & Subrouters

```go
// Restrict the request handler to specific path prefixes.
bookrouter := r.PathPrefix("/books").Subrouter()
bookrouter.HandleFunc("/", AllBooks)
bookrouter.HandleFunc("/{title}", GetBook)
```

## Code

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        title := vars["title"]
        page := vars["page"]

        fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
    })

    http.ListenAndServe(":80", r)
}
```
