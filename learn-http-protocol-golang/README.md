This is ThePrimeagen's course to implement the http protocol in GO.

[Course link](https://www.boot.dev/courses/learn-http-protocol-golang)

Sub sections

[CH4 - L2](https://www.boot.dev/lessons/d35889e3-41d6-4cab-9d93-c30c1d608cc7)

- [x] Remove your useless assert from the request_test.go file.
- [x] Add the following structs to request.go:

```go
// Test: Good GET Request line
r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Good GET Request line with path
r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Invalid number of parts in request line
_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.Error(t, err)
```

- [x] Implement the RequestFromReader function to parse the request-line from the reader. Here are some things to keep in mind:
  - [x] For now, you can slurp the entire request into memory using io.ReadAll and work with the entire thing as a string.
  - [x] Create a parseRequestLine function to do the parsing.
  - [x] Remember that newlines in HTTP are \r\n, not just \n.
  - [x] You can discard everything that comes after the request-line for now.
  - [x] There are always just 3 parts to the request line: strings.Split is your friend here.
  - [x] Verify that the "method" part only contains capital alphabetic characters.
  - [x] Verify that the http version part is 1.1, extracted from the literal HTTP/1.1 format, as we only support HTTP/1.1 for now.
  - [x] Here are some additional references that might help (but don't overthink it, this step should be pretty straightforward):
    - [x] [2.3](https://datatracker.ietf.org/doc/html/rfc9112#section-2.3)
    - [x] [3.1](https://datatracker.ietf.org/doc/html/rfc9112#name-method)
    - [x] [3.2](https://datatracker.ietf.org/doc/html/rfc9112#name-request-target)
- [x] Add more test cases to request_test.go to cover any edge cases you can think of. Here are the names of all the tests I wrote:
  - [x] Good Request line
  - [x] Good Request line with path
  - [x] Good POST Request with path
  - [x] Invalid number of parts in request line
  - [x] Invalid method (out of order) Request line
  - [x] Invalid version in Request line
