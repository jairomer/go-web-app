# Buiding APIs in Go

[source](https://www.youtube.com/watch?v=tIm8UkSf6RA)

## Buiding Web APIs in Golang

**Rules**
- Stay in the standard library as long as possible to avoid external dependencies.
- Use `http.Handler` and `http.HandlerFunc` interfaces.
- Use `http.ResponseWriter` to construct an HTTP response.

```go
type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}

type handlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r)
func (f HandlerFunc) ServeHTTP (w ResponseWriter, r *Request) {
  f(w, r)
}

type ResponseWriter interface {
  Header() Header
  Write([]byte) (int, error)
  WriteHeader(int)
}
```

**http.Request**: Contains everything you need to know about the request.
- `r.Method` - HTTP method
- `r.URL.Path` - Request path
- `r.URL.String()` - Full URL
- `r.URL.Query()` - Query parameters (q=something&p=2)
- `r.Body` - io.ReadCloser of the request body

**Routing**
- Use the standard libary until you need something more.
- Gotyas
  * `/` matches everything
  * No path parameter parsing
  * One handler for every HTTP method
- You might prefer to use Gorilla Mux instead.

**Responding with data**
- Grow your own respond function, but there is an external package called `matryer/respond`.
  + This package provides a way of setting up some additional Options.
- Mirror ServeHTTP signature, even through we are not using *http.Request yet.

```go
func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
  var buf bytes.Buffer
  if err := json.NewEncoder(&buf).Encode(data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(status)
  if _, err := io.Copy(w, &buf); err != nil {
    log.Println("respond:", err)
  }
}

// Example
func handleSomething(w http.ResponseWriter, r *http.Request) {
  data, err := LoadSomething()
  if err != nil {
    respond(w, r, http.StatusInternalServerError, err)
    return
  }
  respond(w, r, http.StatusOk, data)
}
```

## Public Pattern

```go
type Public interface {
  Public() interface {}
}


func respond(w http.ResponseWriter, r *httpRequest, status int, data interface{}) {
  // The following code is a type assertion
  // A type assertion is an operation applied to the value of the interface,
  // or to extract the values of the interface
  //
  // This is used to check that the dynamic type of its operand will match the asserted type or not.
  // If the T is of concrete type, then the type assertion checks that the given dynamic type is indeed T.
  //
  // So essentially, we are making sure that the given object satisfies an interface.
  if obj, ok := data.(Public); ok {
    data = obj.Public()
  }
  // ...
}

type User struct {
  ID            OurID
  Name          string
  PasswordHash  string
}

func (u *User) Public() interface{} {
  return map[string]interface{}{
    "id": u.ID.Encode(),
    "name": u.Name,
  }
}
```

## Decode data from the request
- In order to do this, it is useful to have a single `decode` function.
- This function can be later improved to:
  + Check the `Content-Type` header to use different decoders.
  + Allow URL parameters to influence how the decoder works.
  + Go a little further and validate the input too...

```go
// decode can be this simple to start with, but can be extended
// later to support different formats and behaviours without
// changing the interface.
//
func decode(r *http.Request, v interface{}) error {
  // Notice that you can implement here a Strategy pattern with different decoders.
  if err := json.NewDecoder(r.Body).Decode(v); err != nil {
    return err
  }
  return nil
}
```

## Ok Pattern Example

```go
type Gopher struct {
  Name string
  Country string
}

// We can use this pattern to implement validators
// over a certain input.
func (g *Gopher) OK() error {
  if len(g.Name) == 0 {
    return ErrRequired("name")
  }
  if len(g.Country) == 0 {
    return ErrRequired("country")
  }
  return nil
}

// We are using our previously defined 'decode' function.
func handleCreateGopher(w http.ResponseWriter, r *http.Request) {
  var g gopher
  if err := decode(r, &g); err != nil {
    respond.With(w, r, http.StatusBadRequest, err)
    return
  }
  respond.With(w, r, http.StatusOK, &g)
}

```

## Different types for new things
- New things are different to existing things.

```go
type NewGopher struct {
  Name            string
  Email           string
  Password        string
  PasswordConfirm string
}

type Gopher struct {
  ID              string `json: "id"`
  Name            string `json:"name"`
  Email           string `json:"email"`
  PasswordHash    string `json:"-"`
}

// Save saves a NewGpher and returns the Gopher
func (g *NewGopher) Save(db *mgo.Database) (*Gopher, error) {
  // ...
}
```

## Writing middleware in golang

- A middleware in itself simply takes a `http.HandlerFunc` as one of its parameters, wraps it and returns a new `http.HandlerFunc` for the server to call.
- If implemented correctly, middleware units are extremelly flexible, reusable and sharable.
- Good middleware should not deviate at all from the standard library, the `http.Handler` type is sacred.
  + This allows any Go programmer familiarized with the Handler interface to jump right into your code and be able to reuse it.
- Rather than take in a `http.HandlerFunc` and return a wrapper, we will return a type of this kind.
- Do not break the interface.
Run Code before and after handler code.
```go
func Wrap(h http.Handler) http.Handler {
  return &wrapper{handler: h}
}

type wrapper struct {
  handler http.Handler
}

func (h *wrapper) ServeHTTP(w http.ResponseWritter, r *http.Request) {
  // TODO: do something before each request.
  h.handler.ServeHTTP(w, r)
  // TODO: do something after each request.
}

func main() {
  handler := NewMyHandler()
  http.Handle("/", Wrap(handler))
}
```

### Sharing data across handlers

We can use a Gorilla package called `gorilla/context`:
- Does not break the interface.
- Provides a `map[string]interface{}` per Request
- Easy to clean up (using its own http.Handler wrapper)
- Use `context.Set` and `context.Get` to pass around things.

```go
// We will use the context for the request in order to pass on data.
func prepareSomething(w http.ResponseWriter, r *http.Request) {
  // Store a logger object with the name "logger"
  context.Set(r, "logger", logger)
}

func handleSomethingElse(w http.ResponseWriter, r *http.Request) {
  // Recover a logger object using the name "logger", then assign
  // the interface type.
  logger := context.Get(r, "logger").(*log.Logger)
}
```

**Real World Example for a database connection**
- Almost every API will interact with some kind of datastore.
- Connect when the program is first done (an expensive operation that involves several steps).
- Disconnection when the program is terminated.
- Creation a session per request, which is expected to be relatively cheap.
- Clean up after each request.

```go
func WithDB(s *mgo.Session, h http.Handler) http.Handler {
  return &dbwrapper{dbSession: s, h: h}
}

type dbwrapper struct {
  dbSession *mgo.Session
  h         http.Handler
}

// In order for this to implement the interface 'http.Handler',
// it needs to implement this function.
func (dbwrapper *dbwrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // copy the session
  dbcopy := h.dbSession.Copy()
  defer dbcopy.Close() // clean up after handler execution

  // put session in the context for this Request
  context.Set(r, "db", dbcopy)

  // serve the request
  dbwrapper.h.ServeHTTP(w, r)
}

// Now use it
func main() {
  // connect to mongo
  db, _ := mgo.Dial("localhost") // expensive
  defer db.Close()

  router := mux.NewRouter()
  router.Handle("/things", WithDB(db, http.HandlerFunc(handleThingsRead)))
  router.Handle("/status", http.HandlerFunc(handleStatus))

  http.ListenAndServe(":8080", context.ClearHandler(router))
}

func handleThingsRead(w http.ResponseWriter, r *http.Resquest) {
  db := context.Get(r, "db").(*mgo.Session)

  var results []interface{}
  if err := db.DB("myapp").C("things").Find(nil).All(&results); err != nil {
    respond.With(w, r, http.StatusInternalServerError, err)
    return
  }

  respond.With(w, r, http.StatusOK, results)
}
```

**Request Lifecycle**

1. Client hits `GET \things`
2. `mux.Router` passes the request `WithDB`
3. `WithDB` copies the database session and stored it in the context.
4. `WithDB` then calls the wreapped handler.
5. `handleThingsRead` gets the database session from the context, uses it, and responds.
6. Execution then passes back to `WithFB` which extus and the defferred `Close()` function is called, cleaning up the copy.
7. Finally, `context.ClearHandler` then cleans up the context map for this session.

### Adapter type

```go
type Adapter func(http.Handler) http.Handler
```
Adapter pattern or Decorator pattern.
This is a function that receives a Handler and returns a Handler.

**Example**

```go
func Notify() Adapter {
  return func(h http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
      log.Println("before")
      defer log.Println("after")
      h.ServeHTTP(w, r)
    })
  }
}
```

- This is essentially a Wrapper that returns an Adapter.
- An Adapter is just a function that takes and returns an http.Handler.
  - We can still extend this to specify the type of logger we want to use.

### Using the adapter

Get it and call it.
```go
logger := log.New(os.Stdout, "server:", log.Lshortfile)
http.Handle("/", Notify(loggoer)(indexHandler))
```

Provide yourself an `Adapt` function and use it to take the handler you want to adapt and a List of Adapter types.
It will then iterate over all adapters, calling them one by one in reverse order in a chained manner, returning the result of the first adapter, and finally running all the deferred operations as you return over the stack.
To make the adapters run in the order in which they are specified, you would revers through them in the `Adapt` function.

```go
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
  for _, adapter := range adapters {
    h = adapter(h)
  }
  return h // The nested handler
}

// Usage
//  The execution flow will be such as:
//    - Notify -> Log the event
//    - CopyMgoSession -> Copy the database session and make it available to the handler
//    - CheckAuth -> Check auth credentials and bail if failed.
//    - AddHeader -> Add response header
//    - IndexHandler -> Main business code
//    - Any AddHeader deferred functions
//    - Any CheckAuth deferred functions
//    - CopyMgoSession deferred function (such as closing/commiting the database)
//    - Notify deferred function.
http.Handle("/", Adapt(indexHandler, addHeader("Server", "Mine"),
                                     CheckAuth(providers),
                                     CopyMgoSession(db),
                                     Notify(logger),
                                     ))
```

The nice thing about this pattern is that it does not break the `http.Handler` interface, allowing you to apply middleware when this interface is used and to improve reusability.

## Use Flags For Environment Configuration

```go
func main() {
  var (
    mongoAddr = flag.String("mongo", "", "mongodb addr")
  )
  flag.Parse()
  // ...
}
```

Then run it with: `./myapi -mongo=$DB_PORT_27017_TCP_ADDR`

## Testing

Write unit tests before submitting your code.

### Great Tests

1. Test one thing
2. If that thing breaks, only one test fails.
3. Repeatable tests.
4. Do not rely on run order for tests.
5. Do not require on external resources, tests should be self contained.
  - Mock things or Stub them, then pass them around.
  - Create the interface yourself and wrap an existing external resource.
6. `net/http/httptest` from the standard lib, which is useful for testing handlers.

```go
func TestHandler(t *testing.T) {
  assert := assert.New(t) // testify

  // make test w and r
  r, err := http.NewRequest("GET", "/", nil)
  assert.NoError(err)
  w := httptest.NewRecorder()

  // call handler
  handleSomething(w,r)

  // make assertions
  assert.Equal(w.Body.Strings(), "Hello Golang UK Conference", "body")
  assert.Equal(w.Code, http.StatusOK, "status code")
}

func handleSomething(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "Hello World")
}
```

### Integration testing or table driven testing

Same behaviour as the previous test, but with much less boilerplate.
This technique will allow you to test the full stack, including the middleware and make real HTTP requests.

You will need two components:
- Test table
- Runner

```go

// Test table
var tests = []struct {
    Method        string
    Path          string
    Body          io.Reader
    BodyContains  string
    Status        int
}{{
    Method:       "GET",
    Path:         "/things",
    BodyContains: "Hello Golang UK Conference",
    Status:       http.StatusOK,
}, {
    Method:       "POST",
    Path:         "/things",
    Body:         strings.NewReader({`"name":"Golangg UK Congerence"`}),
    BodyContains: "Hello Golang UK Conference",
    Status:       http.StatusCreated,
}}

// Runner
func TestAll(t *testing.T) {
  assert := assert.New(t)
  server := httptest.NewServer(&myhandler{})
  defer server.Close()
  for _, test := range tests {
    r, err := http.NewRequest(test.Method, server.URL+test.Path, test.Body)
    assert.NoError(err)
    // call handler
    response, err := http.DefaultClient.Do(r)
    assert.NoError(err)
    actualBody, err := ioutil.ReadAll(response.Body)
    assert.NoError(err)
    // make assertions
    assert.Contains(actualBody, test.BodyContains, "body")
    assert.Equal(test.Status, response.StatusCode, "status code")
  }
}
```

### Setup and teardown for each test

```go
var server *httptest.Server

func setup() {
  server = httptest.NewServer(&myAPIHandler{})
}

func teardown() {
  server.Close()
}

func TestSomething(t *testing.T) {
  setup()
  defer teardown()

  // TODO: write test with fresh server
}
```

## Tips and trics

- Expose interfaces, not concrete object.

