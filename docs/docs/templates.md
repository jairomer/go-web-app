# Templates

## Introduction

- Go has a templating language for HTML templates.
- It is mostly used in web applications to display data in a structured way in a client's browser.
- A great benefit of this system is Go's automating escaping of data.
- XXS attacks are mitigated as Go parses the HTML template and escapes all inputs before displaying it to the browser.

## First Template

**Example**
- TODO list as an unordered list in HTML.
- The data passed in can be any kind of Go's data structures.
- It may be a single string or a number or even nested data.
- To access the dat a in a template, use `{{.}}`.
- The dot inside the curly braces is called the pipeline and the root element of the data.

**Example of a nested structure**
```go
data := TodoPageData {
  PageTitle: "My TODO List",
  Todos: []Todo {
    {Title: "task 1", Done: false},
    {Title: "task 2", Done: true},
    {Title: "task 3", Done: true},
  },
}
```

```html
<h1>{{.PageTitle}}</h1>
<ul>
  {{range .Todos}}
    {{if .Done}}
      <li class="done">{{.Title}}</li>
    {{else}}
      <li>{{.Title}}</li>
    {{end}}
  {{end}}
</ul>
```

## Control Structures

- The templatin language contains a righ set of control structured to render your HTML.
- `{{/* This defines a comment */}}`
- `{{.}}` This renders the root element.
- `{{.Title}}` This will render the `Title` field in a nested element.
- `{{if .Done}} {{else}} {{end}}` Defines an if-statement.
- `{{range .Todos}} {{.}} {{end}}` Loops over all `Todos` and renders each using the root element selector.
- `{{block "content" .}} {{end}}` Defines a block with the name "content".

## Parsing Templates from Files

- A template can either be parsed from a string or a file on disk.
- It is usually the case that templates are parses from disk.

**Example**
```go
tmpl, err := template.ParseFiles("layout.html")
// or
tmpl := template.Must(template.ParseFiles("layout.html"))
```

## Execute a Template in a Request Handler
- Once the template is parsed, it is ready to be used in the request handler.
- The `Execute` function accepts an `io.Writer` for writting out the template and an `interface{}` to pass data into the template.
- When the function is called on an `http.ResponseWriter` the Content-Type header is automatically set in the HTTP response to `Content-Type: text/html; charset=utf-8`

```go
func (w http.ResponseWriter, r *http.Request) {
  // ...
  tmpl.Execute(w, "data goes here, bitch")
}
```



