package uuencode

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const (
	// StandardCharset is the standard charset for uuencoded files: ASCII characters 32 - 95.
	StandardCharset = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"

	// AlternateCharset is the same as the standard charset, except that the space character is replaced by backtick.
	// This encoding is non-standard but used occasionally. (Like in the BSD uuencode implementation).
	AlternateCharset = "`!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"
)

// Decoder encapsulates functionality for decoding uuencoded content.
// To create a Decoder, use the helper functions NewStandardDecoder or NewDecoder(charset).
type Decoder struct {
	// encoding is used to decode individual lines within the encoded text.
	encoding *base64.Encoding

	// paddingChar is used to pad lines that have had their padding chopped off for one reason or another.
	paddingChar string
}

// NewStandardDecoder returns a Decoder that uses the StandardCharset.
func NewStandardDecoder() Decoder {
	return NewDecoder(StandardCharset)
}

// NewDecoder returns a decoder using the given charset.
// See StandardCharset and AlternateCharset for common values.
// Note: the provided charset must be a valid base64 charset, otherwise attempts to Decode may panic.
func NewDecoder(charset string) Decoder {
	return Decoder{
		encoding: base64.NewEncoding(charset).WithPadding(base64.NoPadding),
		paddingChar: string(charset[0]), // Padding char is just the first character in the charset
	}
}

// DecodeToBytes is a convenience function for decoding a reader when you just want all the decoded contents in memory in a byte slice.
// See Decode for more info.
func (d Decoder) DecodeToBytes(reader io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	if err := d.Decode(reader, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode decodes the uuencoded contents (as described here: https://en.wikipedia.org/wiki/Uuencoding#Encoded_format)
// of reader and writes the decoded bytes to the given output writer.
// This function assumes there is only one encoded file in the reader, it will ignore anything past the end of the first encoded file.
func (d Decoder) Decode(reader io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(reader)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		if scanner.Err() != nil {
			return fmt.Errorf("error while scanner reader: %w", scanner.Err())
		}

		line := scanner.Text()

		// We don't care about the begin line, we also don't care about empty lines
		if strings.HasPrefix(line, "begin") || len(line) == 0 {
			continue
		}

		// When we find the first end line, we're done.
		if line == "end" {
			return nil
		}

		// Not a special line, let's get to decoding...

		// First check the line length character. If it's the special character backtick (`), the line is empty and we should skip it
		lengthChar := line[0]
		if lengthChar == '`' {
			continue
		}

		// uuencoding adds 32 to the lengthChar so its a printable character
		decodedLen := lengthChar - 32

		// Some encoding schemes don't use the special character for empty lines.
		if decodedLen == 0 {
			continue
		}

		// The formatted characters are everything after the length char.
		// Sometimes padding is omitted from the line, so we have to make sure we add it back before decoding.
		expectedLen := d.encoding.EncodedLen(int(decodedLen))
		encodedCharacters := d.padContentLine(line[1:], expectedLen)

		decoded, err := d.encoding.DecodeString(encodedCharacters)
		if err != nil {
			return fmt.Errorf("error decoding line %d: %w", lineNumber, err)
		}

		// Write the decoded bytes to the output writer
		if _, err := output.Write(decoded[:decodedLen]); err != nil {
			return fmt.Errorf("error writing decoded bytes to writer: %w", err)
		}
	}

	// If we made it out of the loop, it means we never saw the 'end' line
	return fmt.Errorf("malformed input; missing 'end' line")
}

func (d Decoder) padContentLine(line string, expectedLen int) string {
	for len(line) < expectedLen {
		line += d.paddingChar
	}

	return line
}
