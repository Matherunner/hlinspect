package hl

import (
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
)

/*
#include "defs.h"
*/
import "C"

var hlDLL *hooks.Module

var pmPlayerMovePattern = hooks.MakeFunctionPattern("PM_PlayerMove", map[string]string{"Windows": "PM_PlayerMove"}, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 57 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 BA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
	"BigLolly":     hooks.MustMakePattern("55 8B EC 83 EC 0C C7 45 FC 00 00 00 00 A1 ?? ?? ?? ?? 8B 4D 08 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? C7 82 8C 54 04 00 00 00 00 00 A1"),
	"TWHL-Tower-2": hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 8B 4D 08 53 56 57 33 FF 89 7D FC 89 48 04 E8 D8 FC FF FF A1 ?? ?? ?? ?? 89 B8 8C 54 04 00 A1 ?? ?? ?? ?? 0F B6 88 5A 54 04 00"),
})

// HookedPMPlayerMove PM_PlayerMove
//export HookedPMPlayerMove
func HookedPMPlayerMove(server int) {
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

	name, addr, err := pmPlayerMovePattern.Hook(hlDLL, C.HookedPMPlayerMove)
	logs.DLLLog.Debugf("Found %v at %v using %v", pmPlayerMovePattern.Name(), addr, name)

	return
}
