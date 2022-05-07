package game

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/game/cdefs"
	"hlinspect/internal/game/registry"
	"hlinspect/internal/hooks"
	"unsafe"
)

var hlDLL *hooks.Module

// InitHLDLL initialise hl.dll or the corresponding mod DLL
func initHLDLL(base string) (err error) {
	if hlDLL != nil {
		return
	}

	hlDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	reg := Model.Registry()

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&reg.PMInit:                        cdefs.CDefs.HookedPMInit,
		&reg.PMPlayerMove:                  cdefs.CDefs.HookedPMPlayerMove,
		&reg.CSoundEntActiveList:           nil,
		&reg.CSoundEntSoundPointerForIndex: nil,
		&reg.CBaseMonsterChangeSchedule:    nil,
		&reg.CBaseMonsterRouteNew:          nil,
		&reg.WorldGraph:                    nil,
		&reg.CGraphInitGraph:               cdefs.CDefs.CHookedCGraphInitGraph,
		&reg.CBaseMonsterPBestSound:        nil,
	}

	errors := hooks.BatchFind(hlDLL, items)
	printBatchFindErrors(errors)

	switch reg.CBaseMonsterChangeSchedule.PatternKey() {
	case registry.VersionOF8684:
		engine.MonsterOffsets.MonsterState = 0x178
		engine.MonsterOffsets.Schedule = 0x184
		engine.MonsterOffsets.ScheduleIndex = 0x188
		engine.MonsterOffsets.Cine = 0x29c
		engine.MonsterOffsets.AudibleList = 0x228
		engine.MonsterOffsets.WaitFinished = 0x164
		engine.CineOffsets.Radius = 0x2dc
		engine.CineOffsets.Interruptible = 0x2f8
	case registry.VersionHLWON:
		engine.MonsterOffsets.MonsterState = 0x138
		engine.MonsterOffsets.Schedule = 0x144
		engine.MonsterOffsets.ScheduleIndex = 0x148
		engine.MonsterOffsets.Cine = 0x25c
		engine.MonsterOffsets.AudibleList = 0x1e8
		engine.MonsterOffsets.WaitFinished = 0x124
		engine.CineOffsets.Radius = 0x274
		engine.CineOffsets.Interruptible = 0x290
	case registry.VersionOFWON:
		engine.MonsterOffsets.MonsterState = 0x13c
		engine.MonsterOffsets.Schedule = 0x148
		engine.MonsterOffsets.ScheduleIndex = 0x14c
		engine.MonsterOffsets.Cine = 0x260
		engine.MonsterOffsets.AudibleList = 0x1ec
		engine.MonsterOffsets.WaitFinished = 0x128
		engine.CineOffsets.Radius = 0x29c
		engine.CineOffsets.Interruptible = 0x2b8
	case registry.VersionCSCZDS:
		engine.MonsterOffsets.MonsterState = 0x168
		engine.MonsterOffsets.Schedule = 0x174
		engine.MonsterOffsets.ScheduleIndex = 0x178
		engine.MonsterOffsets.Cine = 0x2b0
		engine.MonsterOffsets.AudibleList = 0x23c
		engine.MonsterOffsets.WaitFinished = 0x150
		engine.CineOffsets.Radius = 0x350
		engine.CineOffsets.Interruptible = 0x36c
	case registry.VersionGunman:
		engine.MonsterOffsets.MonsterState = 0x140
		engine.MonsterOffsets.Schedule = 0x14c
		engine.MonsterOffsets.ScheduleIndex = 0x150
		engine.MonsterOffsets.Cine = 0x264
		engine.MonsterOffsets.AudibleList = 0x1f0
		engine.MonsterOffsets.WaitFinished = 0x128
		engine.CineOffsets.Radius = 0x284
		engine.CineOffsets.Interruptible = 0x2a4
	}

	switch reg.CBaseMonsterChangeSchedule.SymbolKey() {
	case registry.VersionWindowsHLDLL:
		engine.CineOffsets.Radius = 0x2a8
		engine.CineOffsets.Interruptible = 0x2c4
	}

	if reg.WorldGraph.Ptr() != nil {
		engine.WorldGraph.SetPointer(reg.WorldGraph.Ptr())
	}

	switch reg.CBaseMonsterRouteNew.PatternKey() {
	case registry.VersionHL8684, registry.VersionCSCZDS, registry.VersionGunman:
		engine.MonsterOffsets.Route = *(*uintptr)(unsafe.Add(reg.CBaseMonsterRouteNew.Ptr(), 0x4)) - 0xc
		engine.MonsterOffsets.RouteIndex = *(*uintptr)(unsafe.Add(reg.CBaseMonsterRouteNew.Ptr(), 0xa))
	}

	return nil
}
