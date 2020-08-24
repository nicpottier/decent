package types

import (
	"math"

	"github.com/nicpottier/decent/parser"
)

/*
ShotSample represents the current state of a shot (or the machine at idle)

  U16    SampleTime       // Time since start of shot, in halfcycles
  U16P12 GroupPressure;   // Pressure at group
  U16P12 GroupFlow;       // Estimated Flow at group
  U16P8  MixTemp;         // Water Temperature entering group
  U24P16 HeadTemp;        // Temperature of water at showerhead
  U16P8  SetMixTemp;      // Set temperature. 0 if no target.
  U16P8  SetHeadTemp;     // Set shower head temp. 0 if no target
  U8P4   SetGroupPressure // Set pressure. 0 if not set.
  U8P4   SetGroupFlow;    // Set flow. 0 if not set.
  U8P0   FrameNumber;     // Frame we are currently in.
  U8P0   SteamTemp;       // Steam metal temp
*/
type ShotSample struct {
	Type             string  `json:"type"`
	SampleTime       int     `json:"sample_time"`
	GroupPressure    float32 `json:"group_pressure"`
	GroupFlow        float32 `json:"group_flow"`
	MixTemp          float32 `json:"mix_temp"`
	HeadTemp         float32 `json:"head_temp"`
	SetMixTemp       float32 `json:"set_mix_temp"`
	SetHeadTemp      float32 `json:"set_head_temp"`
	SetGroupPressure float32 `json:"set_group_pressure"`
	SetGroupFlow     float32 `json:"set_group_flow"`
	FrameNumber      int     `json:"frame_number"`
	SteamTemp        int     `json:"steam_temp"`
}

// ReadShotSample reads and returns a ShotSample from the passed in reader
func ReadShotSample(r *parser.Reader) (interface{}, error) {
	return ShotSample{
		Type:             "shot_sample",
		SampleTime:       int(math.Round(100 * (float64(r.ReadU16()) / (50 * 2)))),
		GroupPressure:    r.ReadU16P12(),
		GroupFlow:        r.ReadU16P12(),
		MixTemp:          r.ReadU16P8(),
		HeadTemp:         r.ReadU24P16(),
		SetMixTemp:       r.ReadU16P8(),
		SetHeadTemp:      r.ReadU16P8(),
		SetGroupPressure: r.ReadU8P4(),
		SetGroupFlow:     r.ReadU8P4(),
		FrameNumber:      r.ReadU8(),
		SteamTemp:        r.ReadU8(),
	}, r.Error()
}
