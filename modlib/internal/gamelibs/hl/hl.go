package hl

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/feed"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"hlinspect/internal/proto"
	"unsafe"
)

/*
#include "defs.h"
*/
import "C"

var hlDLL *hooks.Module

var pmInitPattern = hooks.MakeFunctionPattern("PM_Init", map[string]string{"Windows": "PM_Init"}, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 55 08 33 C0 56 8D 8A ?? ?? ?? ?? 8B B0 ?? ?? ?? ?? 83 C0 0C 89 71 FC 8B B0 ?? ?? ?? ?? 89 31"),
	gamelibs.OF8684: hooks.MustMakePattern("8B 44 24 04 A3 ?? ?? ?? ?? E8 ?? ?? ?? ?? E8 ?? ?? ?? ?? C7 05"),
})
var pmPlayerMovePattern = hooks.MakeFunctionPattern("PM_PlayerMove", map[string]string{"Windows": "PM_PlayerMove"}, map[string]hooks.SearchPattern{
	gamelibs.HL8684:     hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 57 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 BA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
	gamelibs.BigLolly:   hooks.MustMakePattern("55 8B EC 83 EC 0C C7 45 FC 00 00 00 00 A1 ?? ?? ?? ?? 8B 4D 08 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? C7 82 8C 54 04 00 00 00 00 00 A1"),
	gamelibs.TWHLTower2: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 8B 4D 08 53 56 57 33 FF 89 7D FC 89 48 04 E8 D8 FC FF FF A1 ?? ?? ?? ?? 89 B8 8C 54 04 00 A1 ?? ?? ?? ?? 0F B6 88 5A 54 04 00"),
})
var csoundEntActiveListPattern = hooks.MakeFunctionPattern("CSoundEnt::ActiveList", map[string]string{gamelibs.WindowsHLDLL: "CSoundEnt::ActiveList"}, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 58 C3"),
	gamelibs.OF8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 64 C3"),
})
var csoundEntSoundPointerForIndexPattern = hooks.MakeFunctionPattern("CSoundEnt::SoundPointerForIndex", map[string]string{gamelibs.WindowsHLDLL: "CSoundEnt::SoundPointerForIndex"}, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("8B 0D ?? ?? ?? ?? 85 C9 75 03 33 C0 C3 8B 44 24 04 83 F8 3F 7E 13 68 ?? ?? ?? ?? 6A 01 FF 15 ?? ?? ?? ??"),
})
var cbaseMonsterChangeSchedulePattern = hooks.MakeFunctionPattern("CBaseMonster::ChangeSchedule", map[string]string{gamelibs.WindowsHLDLL: "CBaseMonster::ChangeSchedule"}, map[string]hooks.SearchPattern{
	// Search for COND_HEAR_SOUND
	gamelibs.HL8684: hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 78 01 00 00 89 91 7C 01 00 00 89 91 74 01 00 00 89 91 F0 00 00 00 89 91 68 02 00 00"),
	gamelibs.OF8684: hooks.MustMakePattern("8B 81 84 01 00 00 33 D2 3B C2 56 74 55 8B 00 3B C2 74 4F 8B B1 88 01 00 00 57 8B 3C F0"),
})
var cbaseMonsterRouteNewPattern = hooks.MakeFunctionPattern("CBaseMonster::RouteNew", map[string]string{
	gamelibs.WindowsHLDLL: "CBaseMonster::RouteNew",
}, map[string]hooks.SearchPattern{
	// Search for "No Path from %d to %d!" to find CBaseMonster::FGetNodeRoute
	// Go to any of the cross reference, the first function should be CBaseMonster::RouteNew
	// This pattern includes the initial part of CBaseMonster::FRouteClear, and mask out the offsets
	gamelibs.HL8684: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? C1 E0 04"),
})
var worldGraphPattern = hooks.MakeFunctionPattern("WorldGraph", map[string]string{
	// Not actually a function
	gamelibs.WindowsHLDLL: "WorldGraph",
}, nil)
var cgraphInitGraphPattern = hooks.MakeFunctionPattern("CGraph::InitGraph", nil, map[string]hooks.SearchPattern{
	// Search for "Couldn't malloc %d nodes!" to find CGraph::AllocNodes
	// Then find cross reference from CWorld::Precache
	gamelibs.HL8684: hooks.MustMakePattern("56 8B F1 57 33 FF 8B 46 10 89 3E 3B C7 89 7E 04 89 7E 08 74 0C 50 E8 ?? ?? ?? ?? 83 C4 04 89 7E 10 8B 46 0C"),
})

// HookedPMInit PM_Init
//export HookedPMInit
func HookedPMInit(ppm uintptr) {
	hooks.CallFuncInts1(pmInitPattern.Address(), ppm)
	engine.Engine.SetPPMove(unsafe.Pointer(ppm))
	logs.DLLLog.Debugf("Set PPMOVE with address = %x", ppm)
}

// HookedPMPlayerMove PM_PlayerMove
//export HookedPMPlayerMove
func HookedPMPlayerMove(server int) {
	binary, err := proto.Serialize(&proto.PMove{
		Stage:        proto.PMoveStagePre,
		Velocity:     engine.Engine.PMoveVelocity(),
		Position:     engine.Engine.PMovePosition(),
		Viewangles:   engine.Engine.PMoveViewangles(),
		Basevelocity: engine.Engine.PMoveBasevelocity(),
		FSU:          engine.Engine.PMoveCmdFSU(),
		Punchangles:  engine.Engine.PMovePunchangles(),
		EntFriction:  engine.Engine.PMoveEntFriction(),
		EntGravity:   engine.Engine.PMoveEntGravity(),
		FrameTime:    engine.Engine.PMoveFrameTime(),
		Buttons:      engine.Engine.PMoveCmdButtons(),
		Onground:     engine.Engine.PMoveOnground(),
		Flags:        engine.Engine.PMoveFlags(),
		Waterlevel:   engine.Engine.PMoveWaterlevel(),
		InDuck:       engine.Engine.PMoveInDuck(),
		Impulse:      engine.Engine.PMoveImpulse(),
	})
	if err == nil {
		feed.Broadcast(binary)
	}

	hooks.CallFuncInts1(pmPlayerMovePattern.Address(), uintptr(server))
}

// CSoundEntActiveList calls CSoundEnt::ActiveList
func CSoundEntActiveList() int32 {
	return int32(hooks.CallFuncInts0(csoundEntActiveListPattern.Address()))
}

// CSoundEntSoundPointerForIndex calls CSoundEnt::SoundPointerForIndex
func CSoundEntSoundPointerForIndex(index int32) uintptr {
	return uintptr(hooks.CallFuncInts1(csoundEntSoundPointerForIndexPattern.Address(), uintptr(index)))
}

// HookedCGraphInitGraph hooks CGraph::InitGraph
//export HookedCGraphInitGraph
func HookedCGraphInitGraph(this uintptr) {
	hooks.CallFuncThisInts0(cgraphInitGraphPattern.Address(), this)
	engine.WorldGraph.SetPointer(unsafe.Pointer(this))
}

// InitHLDLL initialise hl.dll or the corresponding mod DLL
func InitHLDLL(base string) (err error) {
	if hlDLL != nil {
		return
	}

	hlDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&pmInitPattern:                        C.HookedPMInit,
		&pmPlayerMovePattern:                  C.HookedPMPlayerMove,
		&csoundEntActiveListPattern:           nil,
		&csoundEntSoundPointerForIndexPattern: nil,
		&cbaseMonsterChangeSchedulePattern:    nil,
		&cbaseMonsterRouteNewPattern:          nil,
		&worldGraphPattern:                    nil,
		&cgraphInitGraphPattern:               C.C_HookedCGraphInitGraph,
	}

	errors := hooks.BatchFind(hlDLL, items)
	gamelibs.PrintBatchFindErrors(errors)

	switch cbaseMonsterChangeSchedulePattern.PatternKey() {
	case gamelibs.OF8684:
		engine.MonsterOffsets.Schedule = 0x184
		engine.MonsterOffsets.ScheduleIndex = 0x188
		engine.MonsterOffsets.Cine = 0x29c
		engine.CineOffsets.Radius = 0x2dc
		engine.CineOffsets.Interruptible = 0x2f8
	}

	switch cbaseMonsterChangeSchedulePattern.SymbolKey() {
	case gamelibs.WindowsHLDLL:
		engine.CineOffsets.Radius = 0x2a8
		engine.CineOffsets.Interruptible = 0x2c4
	}

	if worldGraphPattern.Address() != nil {
		engine.WorldGraph.SetPointer(worldGraphPattern.Address())
	}

	switch cbaseMonsterRouteNewPattern.PatternKey() {
	case gamelibs.HL8684:
		engine.MonsterOffsets.Route = *(*uintptr)(unsafe.Pointer(uintptr(cbaseMonsterRouteNewPattern.Address()) + 0x4)) - 0xc
		engine.MonsterOffsets.RouteIndex = *(*uintptr)(unsafe.Pointer(uintptr(cbaseMonsterRouteNewPattern.Address()) + 0xa))
	}

	return
}
