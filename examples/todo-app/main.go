package main

import (
  "html/template"
  "net/http"
  "fmt"
)

type Todo struct {
  Title string
  Done bool
}

type TodoPageData struct {
  PageTitle string
  Todos []Todo
}

func main() {
  tmpl := template.Must(template.ParseFiles("layout.html"))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data := TodoPageData {
      PageTitle: "My TODO List",
      Todos: []Todo {
        {Title: "Task 1", Done: false},
        {Title: "Task 2", Done: true},
        {Title: "Task 3", Done: true},
      },
    }
    // tmpl variable is captured on this scope.
    tmpl.Execute(w, data)
  })
  err := http.ListenAndServe(":8081", nil)
  if err != nil {
    fmt.Println(err)
  }
}
