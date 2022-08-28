package main

/*
A middleware in itself simply takes a http.HandlerFunc as one of its
parameters, wraps it and returns a new http.HandlerFunc for the server to call.

The following code implements a pattern used to build APIs in Golang.

https://www.youtube.com/watch?v=tIm8UkSf6RA
https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81


*/
