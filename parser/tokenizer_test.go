package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tcs := []struct {
		Input   string
		Decoded []byte
		Error   string
	}{
		{"0B6F", []byte{0x0b, 0x6f}, ""},
		{"0B6", nil, "encoding/hex: odd length hex string"},
		{"0N6U", nil, "encoding/hex: invalid byte: U+004E 'N'"},
	}

	for i, tc := range tcs {
		decoded, err := decodeHex([]byte(tc.Input))
		assert.Equal(t, decoded, tc.Decoded, "%d: mismatch on decode", i)

		if tc.Error != "" {
			assert.Equal(t, err.Error(), tc.Error, "%d: mismatch error")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestTokenizer(t *testing.T) {
	tcs := []struct {
		Input  string
		Types  []string
		Bodies [][]byte
		Errors []string
	}{
		{
			"[Q]157AF0",
			[]string{"Q"},
			[][]byte{{0x15, 0x7a, 0xf0}},
			nil,
		},
		{
			"[]",
			nil,
			nil,
			[]string{"invalid message, too short: []"},
		},
		{
			"(a)",
			nil,
			nil,
			[]string{"invalid message type prefix: (a)"},
		},
		{
			"arst\n[Q]0A01",
			[]string{"Q"},
			[][]byte{{0x0a, 0x01}},
			[]string{"invalid message type prefix: arst"},
		},
		{
			"arst\n[Q]0A01\n[P]0BF0",
			[]string{"Q", "P"},
			[][]byte{{0x0a, 0x01}, {0x0b, 0xf0}},
			[]string{"invalid message type prefix: arst"},
		},
	}

	for ti, tc := range tcs {
		reader := strings.NewReader(tc.Input)

		types := make([]string, 0, len(tc.Types))
		bodies := make([][]byte, 0, len(tc.Bodies))
		errors := make([]string, 0, len(tc.Errors))

		var err error
		var m string
		var mb []byte
		var mt string
		for err != io.EOF {
			m, err = ReadNextToken(reader)
			if err == io.EOF {
				continue
			}

			mt, mb, err = splitMessage(m)

			if mt != "" {
				types = append(types, mt)
			}

			if mb != nil {
				bodies = append(bodies, mb)
			}

			if err != nil && err != io.EOF {
				errors = append(errors, err.Error())
			}
		}

		if len(tc.Bodies) == len(bodies) {
			for i := range tc.Bodies {
				assert.Equal(t, tc.Types[i], types[i], "%d/%d: type mismatch", ti, i)
				assert.Equal(t, tc.Bodies[i], bodies[i], "%d/%d: value mismatch", ti, i)
			}
		} else {
			assert.Equal(t, len(tc.Bodies), len(bodies), "%d: mismatched number of bodies: %v != %v", ti, tc.Bodies, bodies)
			assert.Equal(t, len(tc.Types), len(types), "%d: mismatched number of types: %v != %v", ti, tc.Types, types)
		}

		if len(tc.Errors) == len(errors) {
			for i := range tc.Errors {
				assert.Equal(t, tc.Errors[i], errors[i], "%d/%d: error mismatch", ti, i)
			}
		} else {
			assert.Equal(t, len(tc.Errors), len(errors), "%d: mismatched number of errors: %v != %v", ti, tc.Errors, errors)
		}
	}
}
