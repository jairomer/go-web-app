# Assets and files

Now we have a REST endpoint up and running. But we might also want to serve statif files like CSS, javascript or images.

**static-files.go**
```go
package main

import "net/http"

func main() {
  // Initialize file server
  fs := http.FileServer(http.Dir("assets/"))
  // Setup a handle for static files.
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  // Initialize go http server.
  http.ListenAndServe(":8080", nil)
}
```

Now you will be able to fetch static resources as:
```bash
curl -s http://localhost:8080/static/css/styles.css
```
