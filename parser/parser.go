package parser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// ParseMessage parses a message from the passed in bytes according to the passed in type
func ParseMessage(mt string, mb []byte) (interface{}, error) {
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
