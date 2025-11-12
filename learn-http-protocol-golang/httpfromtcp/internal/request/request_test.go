package request

import (
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
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
	})

	t.Run("Good GET Request line with path", func(t *testing.T) {
		t.Parallel()
		r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	})

	t.Run("Invalid number of parts in request line", func(t *testing.T) {
		t.Parallel()
		_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
		require.Error(t, err)
	})
}
