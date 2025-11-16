package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	t.Run("Valid single header", func(t *testing.T) {
		headers := NewHeaders()
		line := "Host: localhost:42069"
		err := headers.ParseLine(line)
		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["Host"])
	})

	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		line := "       Host : localhost:42069       "
		err := headers.ParseLine(line)
		require.Error(t, err)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("Invalid characters in field name", func(t *testing.T) {
		headers := NewHeaders()
		line := "H@st: loclahost:42069"
		err := headers.ParseLine(line)
		require.Error(t, err)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("Invalid delimiter in field name", func(t *testing.T) {
		headers := NewHeaders()
		line := "H,st: loclahost:42069"
		err := headers.ParseLine(line)
		require.Error(t, err)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("Converts field names to canonical", func(t *testing.T) {
		headers := NewHeaders()
		line := "content-type: text/html"
		err := headers.ParseLine(line)
		require.NoError(t, err)
		assert.Equal(t, "text/html", headers["Content-Type"])
	})

	t.Run("Combines the same headers together", func(t *testing.T) {
		headers := NewHeaders()
		err1 := headers.ParseLine("Set-Person: lane-loves-go")
		err2 := headers.ParseLine("Set-Person: prime-loves-zig")
		err3 := headers.ParseLine("Set-Person: tj-loves-ocaml")
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["Set-Person"])
	})
}
