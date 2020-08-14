package types

import (
	"fmt"

	"github.com/nicpottier/decent/parser"
)

// States are the valid states the machine can be in
var States = map[int]string{
	0x00: "sleep",           // 0 Everything is off
	0x01: "going_to_sleep",  // 1 Going to sleep
	0x02: "idle",            // 2 Heaters are controlled, tank water will be heated if required.
	0x03: "busy",            // 3 Firmware is doing something you can't interrupt (eg. cooling down water heater after a shot, calibrating sensors on startup).
	0x04: "espresso",        // 4 Making espresso
	0x05: "steam",           // 5 Making steam
	0x06: "hot_water",       // 6 Making hot water
	0x07: "short_cal",       // 7 Running a short calibration
	0x08: "self_test",       // 8 Checking as much as possible within the firmware. Probably only used during manufacture or repair.
	0x09: "long_cal",        // 9 Long and involved calibration, possibly involving user interaction. (See substates below, for cases like that).
	0x0a: "descale",         // A Descale the whole bang-tooty
	0x0b: "fatal_error",     // B Something has gone horribly wrong
	0x0c: "init",            // C Machine has not been run yet
	0x0d: "no_request",      // D State for T_RequestedState. Means nothing is specifically requested
	0x0e: "skip_to_next",    // E In Espresso, skip to next frame. Others, go to Idle if possible
	0x0f: "hot_water_rinse", // F Produce hot water at whatever temperature is available
	0x10: "steam_rinse",     // 10 Produce a blast of steam
	0x11: "refill",          // 11 Attempting, or needs, a refill.
	0x12: "clean",           // 12 Clean group head
	0x13: "in_boot_loader",  // 13 The main firmware has not run for some reason. Bootloader is active.
	0x14: "air_purge",       // 14 Air purge.
	0x15: "sched_idle",      // 15 Scheduled wake up idle state
}

// SubStates are the valid substates the machine can be in
var SubStates = map[int]string{
	0x00: "no_state",           // 0 State is not relevant
	0x01: "heat_water_tank",    // 1 Cold water is not hot enough. Heating hot water tank.
	0x02: "heat_water_heater",  // 2 Warm up hot water heater for shot.
	0x03: "stabilize_mix_temp", // 3 Stabilize mix temp and get entire water path up to temperature.
	0x04: "pre_infuse",         // 4 Espresso only. Hot Water and Steam will skip this state.
	0x05: "pour",               // 5 Not used in Steam
	0x06: "flush",              // 6 Espresso only, atm
	0x07: "steaming",           // 7 Steam only
	0x08: "descale_int",        // 8 Starting descale
	0x09: "descale_fill_group", // 9 get some descaling solution into the group and let it sit
	0x0a: "descale_return",     // A descaling internals
	0x0b: "descale_group",      // B descaling group
	0x0c: "descale_steam",      // C descaling steam
	0x0d: "clean_init",         // D Starting clean
	0x0e: "clean_fill_group",   // E Fill the group
	0x0f: "clean_soak",         // F Wait for 60 seconds so we soak the group head
	0x10: "clean_group",        // 10 Flush through group
	0x11: "paused_refill",      // 11 Have we given up on a refill
	0x12: "paused_steam",       // 12 Are we paused in steam?

	200: "error_nan",           // 200 Something died with a NaN
	201: "error_inf",           // 201 Something died with an Inf
	202: "error_generic",       // 202 An error for which we have no more specific description
	203: "error_acc",           // 203 ACC not responding, unlocked, or incorrectly programmed
	204: "error_tsensor",       // 204 We are getting an error that is probably a broken temperature sensor
	205: "error_psensor",       // 205 Pressure sensor error
	206: "error_wlevel",        // 206 Water level sensor error
	207: "error_dip",           // 207 DIP switches told us to wait in the error state.
	208: "error_assertion",     // 208 Assertion failed
	209: "error_unsafe",        // 209 Unsafe value assigned to variable
	210: "error_invalid_param", // 210 Invalid parameter passed to function
	211: "error_flash",         // 211 Error accessing external flash
	212: "error_oom",           // 212 Could not allocate memory
	213: "error_deadline",      // 213 Realtime deadline missed
}

/*
StateInfo represents the state and substate the machine is currently in.
   U8 State     // The state the machine is currently in
   U8 SubState  // The substate the machine is currently in
*/
type StateInfo struct {
	Type     string `json:"type"`
	State    string `json:"state"`
	SubState string `json:"substate"`
}

// ReadStateInfo reads and returns StateInfo
func ReadStateInfo(r *parser.Reader) (interface{}, error) {
	s := r.ReadU8()
	state, found := States[s]
	if !found {
		state = fmt.Sprintf("unknown state: %d", s)
	}

	s = r.ReadU8()
	substate, found := SubStates[s]
	if !found {
		state = fmt.Sprintf("unknown substate: %d", s)
	}

	return StateInfo{
		Type:     "state_info",
		State:    state,
		SubState: substate,
	}, r.Error()
}
