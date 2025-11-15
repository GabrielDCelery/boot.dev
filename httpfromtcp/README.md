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

[CH4 - L3](https://www.boot.dev/lessons/daf467c5-f17b-4382-a74c-e4bc68f2fc8d)

- [x] Paste this code into your request_test.go file:

```go
type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}
```

- [x] Update your test suite to use the chunkReader type and test for different numbers of bytes read per chunk:

```go
// Test: Good GET Request line
reader := &chunkReader{
	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 3,
}
r, err := RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Good GET Request line with path
reader = &chunkReader{
	data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	numBytesPerRead: 1,
}
r, err = RequestFromReader(reader)
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
```

> [!WARN]
> Be sure to test values as low as 1 and as high as the length of the request string. Our code should work under all conditions.

- [x] Update your parseRequestLine to return the number of bytes it consumed. If it can't find an \r\n (this is important!) it should return 0 and no error. This just means that it needs more data before it can parse the request line.
- [x] Add a new internal "enum" (I just used an int) to your Request struct to track the state of the parser. For now, you just need 2 states:
  - "initialized"
  - "done".
- [x] Implement a new func (r \*Request) parse(data []byte) (int, error) method.
  - [x] It accepts the next slice of bytes that needs to be parsed into the Request struct
  - [x] It updates the "state" of the parser, and the parsed RequestLine field.
  - [x] It returns the number of bytes it consumed (meaning successfully parsed) and an error if it encountered one.

[CH4 - L6](https://www.boot.dev/lessons/56df6098-0175-4e83-a481-a5381db3d9fd)

- [] Delete your getLinesChannel function - we're working with HTTP now, not just newlines. Similarly, delete the resulting logic that prints the text coming back across its channel.
- [] Instead, call RequestFromReader. Assuming it's successful, print out the RequestLine in this format (with the dynamic data of course):

```sh
Request line:
- Method: GET
- Target: /
- Version: 1.1
```

- [] Run your tcplistener program again from the root of your project and redirect the output to a temporary file:

```sh
go run ./cmd/tcplistener | tee /tmp/requestline.txt
```

- [] In another shell, send this request to it:

```sh
curl http://localhost:42069/prime/agen
```

- [] Kill both programs. Your requestline.txt file should contain the output of your tcplistener program.
