# HTTP Server

A basic HTTP server should be able to:
- *Process Dynamic Requests* from users who browse the website, log into their accounts or post images.
- *Serve static assets* such as JavaScript, CSS and images to browsers to create a dynamic experience for the user.
- *Accept connections* when the HTTP Server is listening to a specific port accessible from the Internet.

## Process Dynamic Requests

In order to do this we accept requests in terms of the resource URL being accessed to using handlers.

```go
http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    // If GET method then
    // r.URL.Query().Get("token")

    // If POST method then
    // r.FormValue("email")

    fmt.Fprint(w, "Welcome to my website!")
})
```

## Serving static assets

- In order to achieve this we will use the inbuilt `http.FileServer` and point it to a URL path.
- For the file server to work, it needs to know where to serve files from.

```go
fs := http.FileServer(http.Dir("static/"))
```

Once our server is in place, we just need to point a url path at it.
In order to serve files correctly, we need to strip away a part of the URL path.
Usually this is the name of the directory our files live in.
```go
http.Handle("/static/", http.StripPrefix("/static/", fs))
```

## Accept connections

Go has an inbuilt http server, so this is straighforward.

```go
http.ListenAndServe(":<port>", nil)
```

## Code

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to my website!")
    })

    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(":80", nil)

```
