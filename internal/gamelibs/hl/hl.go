package hl

import (
	"fmt"
	"hlinspect/internal/engine"
	"hlinspect/internal/feed"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"unsafe"
)

/*
#include "defs.h"
*/
import "C"

var hlDLL *hooks.Module

var pmInitPattern = hooks.MakeFunctionPattern("PM_Init", map[string]string{"Windows": "PM_Init"}, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 55 08 33 C0 56 8D 8A ?? ?? ?? ?? 8B B0 ?? ?? ?? ?? 83 C0 0C 89 71 FC 8B B0 ?? ?? ?? ?? 89 31"),
})
var pmPlayerMovePattern = hooks.MakeFunctionPattern("PM_PlayerMove", map[string]string{"Windows": "PM_PlayerMove"}, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 57 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 BA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
	"BigLolly":     hooks.MustMakePattern("55 8B EC 83 EC 0C C7 45 FC 00 00 00 00 A1 ?? ?? ?? ?? 8B 4D 08 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? C7 82 8C 54 04 00 00 00 00 00 A1"),
	"TWHL-Tower-2": hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 8B 4D 08 53 56 57 33 FF 89 7D FC 89 48 04 E8 D8 FC FF FF A1 ?? ?? ?? ?? 89 B8 8C 54 04 00 A1 ?? ?? ?? ?? 0F B6 88 5A 54 04 00"),
})

// HookedPMInit PM_Init
//export HookedPMInit
func HookedPMInit(ppm uintptr) {
	hooks.CallFuncInts1(pmInitPattern.Address(), ppm)
	engine.Engine.SetPPMove(ppm)
	logs.DLLLog.Debugf("Set PPMOVE with address = %x", ppm)
}

// HookedPMPlayerMove PM_PlayerMove
//export HookedPMPlayerMove
func HookedPMPlayerMove(server int) {
	vel := engine.Engine.PMoveVelocity()
	feed.Broadcast(fmt.Sprintf("%v %v %v", vel[0], vel[1], vel[2]))
	hooks.CallFuncInts1(pmPlayerMovePattern.Address(), uintptr(server))
}

// InitHLDLL initialise hl.dll or the corresponding mod DLL
func InitHLDLL() (err error) {
	if hlDLL != nil {
		return
	}

	hlDLL, err = hooks.NewModule("hl.dll")
	if err != nil {
		return
	}

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&pmInitPattern:       C.HookedPMInit,
		&pmPlayerMovePattern: C.HookedPMPlayerMove,
	}

	errors := hooks.BatchFind(hlDLL, items)
	for pat, err := range errors {
		if err == nil {
			logs.DLLLog.Debugf("Found %v at %v", pat.Name(), pat.Address())
		} else {
			logs.DLLLog.Debugf("Failed to find %v: %v", pat.Name(), err)
		}
	}

	return
}
