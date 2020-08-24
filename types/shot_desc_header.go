package types

import "github.com/nicpottier/decent/parser"

// U8P0 HeaderV;           // Set to 1 for this type of shot description
// U8P0 NumberOfFrames;    // Total number of frames.
// U8P0 NumberOfPreinfuseFrames; // Number of frames that are preinfusion
// U8P4 MinimumPressure;   // In flow priority modes, this is the minimum pressure we'll allow
// U8P4 MaximumFlow;       // In pressure priority modes, this is the maximum flow rate we'll allow
type ShotDescHeader struct {
	Type               string  `json:"type"`
	Version            int     `json:"header_version"`
	NumFrames          int     `json:"num_frames"`
	NumPreinfuseFrames int     `json:"num_preinfuse_frames"`
	MinPressure        float32 `json:"min_pressure"`
	MaxFlow            float32 `json:"max_flow"`
}

// ReadShotDescHeader reads and returns a ShotDescHeader from the passed in stream
func ReadShotDescHeader(r *parser.Reader) (interface{}, error) {
	return ShotDescHeader{
		Type:               "shot_desc_header",
		Version:            r.ReadU8(),
		NumFrames:          r.ReadU8(),
		NumPreinfuseFrames: r.ReadU8(),
		MinPressure:        r.ReadU8P4(),
		MaxFlow:            r.ReadU8P4(),
	}, r.Error()
}
