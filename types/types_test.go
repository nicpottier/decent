package types

import (
	"bytes"
	"testing"

	"github.com/nicpottier/decent/parser"
	"github.com/stretchr/testify/assert"
)

func TestReadingTypes(t *testing.T) {
	tcs := []struct {
		Input string
		Value interface{}
		Error string
	}{
		{
			"[Q]0B6F0500",
			WaterLevels{
				"water_levels",
				11.433594,
				5,
			},
			"",
		},
		{
			"[Q]0B710500",
			WaterLevels{
				"water_levels",
				11.441406,
				5,
			},
			"",
		},
		{
			"[Q]0B7105000A",
			nil,
			"unexpected extra input reading message: [Q]0B7105000A",
		},
		{
			"[Q]0B",
			nil,
			"unexpected EOF",
		},
		{
			"[M]9924000A000019F61AE4591B5F590000600325",
			ShotSample{
				Type:             "shot_sample",
				SampleTime:       39204,
				GroupPressure:    0.0024414062,
				GroupFlow:        0,
				MixTemp:          25.960938,
				HeadTemp:         26.891983,
				SetMixTemp:       27.371094,
				SetHeadTemp:      89,
				SetGroupPressure: 0,
				SetGroupFlow:     6,
				FrameNumber:      3,
				SteamTemp:        37,
			},
			"",
		},
		{
			"[M]3A14000400005A2D58CC3250005900003003A4",
			ShotSample{
				Type:             "shot_sample",
				SampleTime:       14868,
				GroupPressure:    0.0009765625,
				GroupFlow:        0,
				MixTemp:          90.17578,
				HeadTemp:         88.79764,
				SetMixTemp:       80,
				SetHeadTemp:      89,
				SetGroupPressure: 0,
				SetGroupFlow:     3,
				FrameNumber:      3,
				SteamTemp:        164,
			},
			"",
		},
		{
			"[N]0402",
			StateInfo{
				Type:     "state_info",
				State:    "espresso",
				SubState: "heat_water_heater",
			},
			"",
		},
	}

	for ti, tc := range tcs {
		ms, err := parser.ReadNextToken(bytes.NewReader([]byte(tc.Input)))
		assert.NoError(t, err)
		m, err := parser.ParseMessage(ms)

		assert.Equal(t, tc.Value, m, "%d: mismatched value", ti)
		if tc.Error != "" {
			assert.Equal(t, err.Error(), tc.Error, "%d: mismatch error")
		} else {
			assert.NoError(t, err)
		}
	}
}
