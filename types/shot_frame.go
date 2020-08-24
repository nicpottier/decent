package types

import (
	"github.com/nicpottier/decent/parser"
)

// ShotFrame describes a single frame in a profile
//	U8P0   Flag       // See T_E_FrameFlags
//	U8P4   SetVal     // SetVal is a 4.4 fixed point number, setting either pressure or flow rate, as per mode
//	U8P1   Temp       // Temperature in 0.5 C steps from 0 - 127.5
//	F8_1_7 FrameLen   // FrameLen is the length of this frame in seconds. It's a 1/7 bit floating point number as described in the F8_1_7 a struct
//	U8P4   TriggerVal // Trigger value. Could be a flow or pressure.
//	U10P0  MaxVol     // Exit current frame if the volume/weight exceeds this value. 0 means ignore
type ShotFrame struct {
	Type        string `json:"type"`
	FrameNumber int    `json:"frame_number"`

	TempType string  `json:"temp_type,omitempty"`
	Temp     float32 `json:"temp"`

	ControlType  string  `json:"control_type"`
	ControlValue float32 `json:"control_value"`

	Trigger         bool    `json:"trigger,omitempty"`
	TriggerType     string  `json:"trigger_type,omitempty"`
	TriggerOperator string  `json:"trigger_operator,omitempty"`
	TriggerValue    float32 `json:"trigger_value,omitempty"`

	Interpolate  bool `json:"interpolate,omitempty"`
	IgnoreLimits bool `json:"ignore_limits,omitempty"`

	MaxVolume int `json:"max_volume,omitempty"`

	FrameLength float32 `json:"frame_length"`

	FlagBits int `json:"_flag_bits"`
}

// bitmasks on our flags
const (
	ControlTypeMask     = 1 << 0 // Are we trying to control pressure or flow?
	TriggerMask         = 1 << 1 // Do a compare, early exit current frame if compare true
	TriggerOperatorMask = 1 << 2 // If we are doing a compare, then 0 = less than, 1 = greater than
	TriggerTypeMask     = 1 << 3 // Compare Pressure or Flow?
	TempTargetMask      = 1 << 4 // Disable shower head temperature compensation. Target Mix Temp instead.
	InterpolateMask     = 1 << 5 // Hard jump to target value, or ramp?
	IgnoreLimitsMask    = 1 << 6 // Ignore minimum pressure and max flow settings
)

var controlTypeValues map[int]string = map[int]string{0: "pressure", ControlTypeMask: "flow"}
var triggerValues map[int]bool = map[int]bool{0: false, TriggerMask: true}
var triggerOperatorValues map[int]string = map[int]string{0: "less_than", TriggerOperatorMask: "greater_than"}
var triggerTypeValues map[int]string = map[int]string{0: "pressure", TriggerTypeMask: "flow"}
var tempTypeValues map[int]string = map[int]string{0: "basket", TempTargetMask: "mix"}
var interpolateValues map[int]bool = map[int]bool{0: false, InterpolateMask: true}
var ignoreLimitsValues map[int]bool = map[int]bool{0: false, IgnoreLimitsMask: true}

func (f *ShotFrame) parseFlags() {
	b := f.FlagBits
	f.ControlType = controlTypeValues[b&ControlTypeMask]
	f.Trigger = triggerValues[b&TriggerMask]
	if f.Trigger {
		f.TriggerType = triggerTypeValues[b&TriggerTypeMask]
		f.TriggerOperator = triggerOperatorValues[b&TriggerOperatorMask]
	}
	f.TempType = tempTypeValues[b&TempTargetMask]
	f.Interpolate = interpolateValues[b&InterpolateMask]
	f.IgnoreLimits = ignoreLimitsValues[b&IgnoreLimitsMask]
}

// ReadShotFrame reads and returns a ShotFramee from the passed in reader
func ReadShotFrame(r *parser.Reader) (interface{}, error) {
	f := ShotFrame{
		Type:         "shot_frame",
		FrameNumber:  r.ReadU8(),
		FlagBits:     r.ReadU8(),
		ControlValue: r.ReadU8P4(),
		Temp:         r.ReadU8P1(),
		FrameLength:  r.ReadF817(),
		TriggerValue: r.ReadU8P4(),
		MaxVolume:    r.ReadU10(),
	}
	f.parseFlags()
	return f, r.Error()
}
