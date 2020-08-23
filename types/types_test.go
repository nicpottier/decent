package types

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/nicpottier/decent/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateTruth = flag.Bool("update", false, "whether to update truth files")

func TestReadingTypes(t *testing.T) {
	flag.Parse()

	type TestCase struct {
		Label   string          `json:"label"`
		Message string          `json:"message"`
		Parsed  json.RawMessage `json:"parsed,omitempty"`
		Error   string          `json:"error,omitempty"`
	}

	var truthFile = "./testdata/type_tests.json"
	tcs := make([]*TestCase, 0, 20)
	tcJSON, err := ioutil.ReadFile(truthFile)
	require.NoError(t, err)

	err = json.Unmarshal(tcJSON, &tcs)
	require.NoError(t, err)

	for ti, tc := range tcs {
		ms, err := parser.ReadNextToken(bytes.NewReader([]byte(tc.Message)))
		assert.NoError(t, err)
		m, err := parser.ParseMessage(ms)

		if *updateTruth {
			if err != nil {
				tc.Error = err.Error()
			} else {
				tc.Error = ""
			}
		} else {
			if tc.Error != "" {
				assert.Equal(t, err.Error(), tc.Error, "%d: mismatch error", ti)
			} else {
				assert.NoError(t, err, "%d: unexpected error", ti)
			}
		}

		if *updateTruth {
			tc.Parsed = nil
			if m != nil {
				tc.Parsed, _ = json.Marshal(m)
			}
		} else {
			if m != nil {
				mj, _ := json.MarshalIndent(m, "", "  ")
				assert.JSONEq(t, string(mj), string(tc.Parsed))
			}
		}
	}

	if *updateTruth {
		truth, err := json.MarshalIndent(tcs, "", "  ")
		require.NoError(t, err)

		err = ioutil.WriteFile(truthFile, truth, 0644)
		require.NoError(t, err, "failed to update truth file")
	}
}
