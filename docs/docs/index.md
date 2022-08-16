# Go Web Examples

## Hello World

- Go is a battery included programming language that has a web server already built in.
- The `net/http` package contains an http client and server.

### Registering a Request Handler

A handler in Go is a function with this signature.
```go
// http.ResponseWriter is where you write the text/html response to
// http.Request contains all the information about this HTTP request.
func (w http.ResponseWriter, r *http.Request)
```

Registering a request handler is as simple as:
```go
http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
})
```

### Listen for HTTP Connections

- The request handler alone cannot accept any HTTP connections from the outside.
- An HTTP server has to listen on a port to pass connections on to the request handler.

```go
http.ListenAndServe(":80", nil)
```

### Complete code

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
    })

    http.ListenAndServe(":80", nil)
}
```
