package parser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// ParseMessage parses a message from the passed in bytes according to the passed in type
func ParseMessage(ms string) (interface{}, error) {
	mt, mb, err := splitMessage(ms)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to split message: %s", ms)
	}

	tr, found := typeMap[mt]
	if !found {
		return nil, fmt.Errorf("unable to find reader for type: %s", mt)
	}

	br := bytes.NewReader(mb)
	r := NewReader(br)
	m, err := tr(r)
	if err != nil {
		return nil, err
	}

	// if we aren't at the end of our buffer, that's also an error, we should consume the entire message
	_, err = br.ReadByte()
	if err != io.EOF {
		return nil, fmt.Errorf("unexpected extra input reading message: [%s]%s", mt, strings.ToUpper(hex.EncodeToString(mb)))
	}

	return m, nil
}

// checks that our type prefix looks valid, strips it out, then returns the type and decoded body
func splitMessage(m string) (string, []byte, error) {
	if len(m) < 3 {
		return "", nil, fmt.Errorf("invalid message, too short: %s", string(m))
	}

	if (m[0] != '[' && m[2] != ']') && (m[0] != '<' && m[2] != '>') {
		return "", nil, fmt.Errorf("invalid message type prefix: %s", string(m))
	}

	body, err := decodeHex([]byte(m[3:]))
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

// RegisterReader allows for type readers to register themselves as readers fo the passed in type
func RegisterReader(t string, f TypeReader) {
	_, found := typeMap[t]
	if found {
		panic(fmt.Sprintf("duplicate type: %s", t))
	}

	typeMap[t] = f
}

// TypeReader defines the interface to read a type
type TypeReader func(*Reader) (interface{}, error)

var typeMap = map[string]TypeReader{}
