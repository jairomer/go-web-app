package main

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
)

func main() {

  r := mux.NewRouter()

  r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you have requested: %s\n", r.URL.Path)
    vars := mux.Vars(r)
    title := vars["title"]
    page := vars["page"]
    fmt.Fprintf(w, "Hello, you have requested the book %s on page %s\n", title, page)
  })

  fs := http.FileServer(http.Dir("static/"))
  r.Handle("/static/", http.StripPrefix("/static/", fs))

  err := http.ListenAndServe(":8081", r)
  if err != nil {
    fmt.Println(err)
  }
}
