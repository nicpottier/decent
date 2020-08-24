package types

import "github.com/nicpottier/decent/parser"

// 0x80: Start the steam quickly and at higher pressure
// 0x00: Start the steam slowly and at lower pressure (ie. No Bit set)
// 0x40: Run the steam at higher pressure
// 0x00: Run the steam at lower pressure
type SteamSettings struct {
	Start string `json:"start"`
	Power string `json:"power"`
}

const (
	startMask = 0x80
	powerMask = 0x40
)

var startValues = map[int]string{0: "slow", startMask: "fast"}
var powerValues = map[int]string{0: "low", powerMask: "high"}

func readSteamSettings(b int) SteamSettings {
	return SteamSettings{
		Start: startValues[b&startMask],
		Power: powerValues[b&powerMask],
	}
}

// T_Enum_SteamSetting SteamSettings; // Defines the steam shot
// U8P0  TargetSteamTemp; // Valid range is 140 - 160
// U8P0  TargetSteamLength;  // Length in seconds of steam
// U8P0  TargetHotWaterTemp; // Temperature of the mixed hot water
// U8P0  TargetHotWaterVol;     // How much water we'll need for hot water (so we know if we have enough)
// U8P0  TargetHotWaterLength;  // (DE1 only) Length of time for a shot (water vol is ignored)
// U8P0  TargetEspressoVol;     // So we know if we have enough water
// U16P8 TargetGroupTemp;       // So we know what to set the group to
type ShotSettings struct {
	Type                 string        `json:"type"`
	SteamSettings        SteamSettings `json:"steam_settings"`
	SteamBits            int           `json:"_steam_bits"`
	TargetSteamTemp      int           `json:"target_steam_temp"`
	TargetSteamLength    int           `json:"target_steam_length"`
	TargetWaterTemp      int           `json:"target_water_temp"`
	TargetWaterVolume    int           `json:"target_water_volume"`
	TargetWaterLength    int           `json:"target_water_length"`
	TargetEspressoVolume int           `json:"target_espresso_volume"`
	TargetGroupTemp      float32       `json:"target_group_temp"`
}

// ReadShotSettings reads and returns a ShotSettings from the passed in stream
func ReadShotSettings(r *parser.Reader) (interface{}, error) {
	s := ShotSettings{
		Type:                 "shot_settings",
		SteamBits:            r.ReadU8(),
		TargetSteamTemp:      r.ReadU8(),
		TargetSteamLength:    r.ReadU8(),
		TargetWaterTemp:      r.ReadU8(),
		TargetWaterVolume:    r.ReadU8(),
		TargetWaterLength:    r.ReadU8(),
		TargetEspressoVolume: r.ReadU8(),
		TargetGroupTemp:      r.ReadU16P8(),
	}
	s.SteamSettings = readSteamSettings(s.SteamBits)
	return s, r.Error()
}
