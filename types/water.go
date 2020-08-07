package types

import "github.com/nicpottier/decent/parser"

/*
WaterLevels represents the current water level states on the machine.

  U16P8 Level;          // Writes to Level will be ignored
  U16P8 StartFillLevel; // Writes to this set the water level refill trigger point.
                        // Invalid writes will be clipped to valid levels
						// To disable refill, set StartFillLevel to 0.0
*/
type WaterLevels struct {
	Type           string  `json:"type"`
	WaterLevel     float32 `json:"water_level"`
	StartFillLevel float32 `json:"water_fill_level"`
}

// ReadWaterLevels reads and returns a WaterLevels message from the passed in reader
func ReadWaterLevels(r *parser.Reader) (interface{}, error) {
	return WaterLevels{
		Type:           "water_levels",
		WaterLevel:     r.ReadU16P8(),
		StartFillLevel: r.ReadU16P8(),
	}, r.Error()
}
