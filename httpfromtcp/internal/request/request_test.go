package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	t.Run("Good GET Request line", func(t *testing.T) {
		t.Parallel()
		r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
	})

	t.Run("Good GET Request line with path", func(t *testing.T) {
		t.Parallel()
		r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	})

	t.Run("Invalid number of parts in request line", func(t *testing.T) {
		t.Parallel()
		_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
	})

	t.Run("Good POST Request Line", func(t *testing.T) {
		t.Parallel()
		r, err := RequestFromReader(strings.NewReader("POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "POST", r.RequestLine.Method)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	})

	t.Run("Invalid method (out of order) Request line", func(t *testing.T) {
		t.Parallel()
		_, err := RequestFromReader(strings.NewReader("/coffee GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
	})

	t.Run("Invalid version in request line", func(t *testing.T) {
		t.Parallel()
		_, err := RequestFromReader(strings.NewReader("GET / HTTP/4\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
	})

	t.Run("Good GET Request line with when reading chunks of 3 bytes", func(t *testing.T) {
		t.Parallel()
		reader := NewChunkReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n", 3)
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("Good GET Request line with when reading chunks of 1 byte", func(t *testing.T) {
		t.Parallel()
		reader := NewChunkReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n", 1)
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	})

	t.Run("Good GET request with headers", func(t *testing.T) {
		t.Parallel()
		reader := &chunkReader{
			data:              "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numOfBytesPerRead: 3,
		}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "localhost:42069", r.Headers["Host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["User-Agent"])
		assert.Equal(t, "*/*", r.Headers["Accept"])
	})

	t.Run("Good GET request with body", func(t *testing.T) {
		t.Parallel()
		reader := &chunkReader{
			data:              "POST /submit HTTP/1.1\r\n" + "Host: localhost:42069\r\n" + "Content-Length: 13\r\n" + "\r\n" + "hello world!\n",
			numOfBytesPerRead: 3,
		}

		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "POST", r.RequestLine.Method)
		assert.Equal(t, "/submit", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "localhost:42069", r.Headers["Host"])
		assert.Equal(t, "13", r.Headers["Content-Length"])
		assert.Equal(t, "hello world!\n", string(r.Body))
	})
}

type chunkReader struct {
	data              string
	numOfBytesPerRead int
	pos               int
}

func NewChunkReader(data string, numOfBytesPerRead int) *chunkReader {
	return &chunkReader{data: data, numOfBytesPerRead: numOfBytesPerRead, pos: 0}
}

func (cr *chunkReader) Read(p []byte) (int, error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	end := min(cr.pos+cr.numOfBytesPerRead, len(cr.data))
	n := copy(p, cr.data[cr.pos:end])
	cr.pos += n
	return n, nil
}
