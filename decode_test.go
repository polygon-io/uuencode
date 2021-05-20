package uuencode

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeToBytes(t *testing.T) {
	t.Run("standard encoding", func(t *testing.T) {
		input := "begin 644 cat.txt\n" +
			"#0V%T\n" +
			"`\n" +
			"end"

		decoder := NewStandardDecoder()

		results, err := decoder.DecodeToBytes(strings.NewReader(input))
		require.NoError(t, err)
		assert.EqualValues(t, "Cat", string(results))
	})

	t.Run("BSD/Alternate style encoding", func(t *testing.T) {
		// polygon.uu was uuencoded using the `uuencode` utility on macOS
		encodedFile, err := os.Open("test_data/polygon.uu")
		require.NoError(t, err)

		defer encodedFile.Close()

		decoder := NewDecoder(AlternateCharset)
		decodedBytes, err := decoder.DecodeToBytes(encodedFile)
		require.NoError(t, err)

		expectedFile, err := os.Open("test_data/polygon.jpg")
		require.NoError(t, err)

		defer expectedFile.Close()

		expectedBytes, err := io.ReadAll(expectedFile)
		require.NoError(t, err)

		assert.Equal(t, expectedBytes, decodedBytes)
	})
}
