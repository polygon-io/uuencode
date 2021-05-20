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
}

func TestDecodeFiles(t *testing.T) {
	t.Run("BSD/Alternate style encoding", func(t *testing.T) {
		// polygon.uu was uuencoded using the `uuencode` utility on macOS
		decodeFile(t, NewDecoder(AlternateCharset), "test_data/polygon.uu", "test_data/polygon.jpg")
	})

	t.Run("standard/SEC style encoding", func(t *testing.T) {
		// polygon.sec.uu was encoded in the same style that the SEC uses when disseminating binary files in filings
		decodeFile(t, NewStandardDecoder(), "test_data/polygon.sec.uu", "test_data/polygon.jpg")
	})
}

func decodeFile(t *testing.T, decoder Decoder, encodedFilename, decodedFilename string) {
	encodedFile, err := os.Open(encodedFilename)
	require.NoError(t, err)

	defer encodedFile.Close()

	decodedBytes, err := decoder.DecodeToBytes(encodedFile)
	require.NoError(t, err)

	expectedFile, err := os.Open(decodedFilename)
	require.NoError(t, err)

	defer expectedFile.Close()

	expectedBytes, err := io.ReadAll(expectedFile)
	require.NoError(t, err)

	assert.Equal(t, expectedBytes, decodedBytes)
}
