package parser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// ReadNextToken reads the next message token in the stream, returning the type and value.
//
// Messages must start with the type surrounded in square brackets. ex: [M], [Q]
//
// Message values follow the type in two character pair hex encodings. ex: 39 -> 0x39
//
// Messages are demarcated by newlines (\n).
//
// Full examples:
// [M]389D000400005A2558CECF50005900003003A4
// [Q]0B790500
//
// This function will block reading the next message off the passed in stream, returning the string
// type and the decoded byte payload for the message.
//
// If a malformed message is read (invalid header, odd length), an error is returned, but the reader
// recovers and can continue to be called until EOF.
//
// End of stream will be indicated by returning io.EOF
func ReadNextToken(reader io.Reader) (string, []byte, error) {
	buffer := bytes.Buffer{}
	b := make([]byte, 1)

	var err error
	var n int

	for {
		n, err = reader.Read(b)

		// oh dear, something went wrong
		if err != nil && err != io.EOF {
			return "", nil, fmt.Errorf("error while reading from buffer: %s", err.Error())
		}

		// end of message, break out
		if b[0] == '\n' {
			break
		}

		// otherwise, add to our buffer and keep reading
		if n > 0 {
			buffer.Write(b)
		}

		// end of stream, break out
		if err == io.EOF {
			// no message read, return immediately with EOF
			if buffer.Len() == 0 {
				return "", nil, io.EOF
			}

			// otherwise break out and parse what we've read
			break
		}

		// nothing read, pause a bit before looping
		if n == 0 {
			time.Sleep(time.Millisecond * 50)
		}
	}

	return splitMessage(buffer.Bytes())
}

// checks that our type prefix looks valid, strips it out, then returns the type and decoded body
func splitMessage(m []byte) (string, []byte, error) {
	if len(m) < 3 {
		return "", nil, fmt.Errorf("invalid message, too short: %s", string(m))
	}

	if m[0] != '[' && m[2] != ']' {
		return "", nil, fmt.Errorf("invalid message type prefix: %s", string(m))
	}

	body, err := decodeHex(m[3:])
	if err != nil {
		return "", nil, fmt.Errorf("error decoding message body: %s", string(m))
	}

	return string(m[1]), body, nil
}

// decodes two pair ascii hex encodings.. 39 -> 0x39  42 -> 0x42
func decodeHex(encoded []byte) ([]byte, error) {
	decoded := make([]byte, len(encoded)/2)

	_, err := hex.Decode(decoded, encoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
