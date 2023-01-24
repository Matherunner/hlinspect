package game

import (
	"hlinspect/internal/game/cdefs"
	"hlinspect/internal/game/engine"
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
		&reg.CBasePlayerForceClientDllUpdate: nil,
	}

	errors := hooks.BatchFind(hlDLL, items)
	printBatchFindErrors(errors)

	switch reg.CBaseMonsterChangeSchedule.PatternKey() {
	case registry.VersionHLWON:
		engine.MonsterOffsets.MonsterState = 0x138
		engine.MonsterOffsets.Schedule = 0x144
		engine.MonsterOffsets.ScheduleIndex = 0x148
		engine.MonsterOffsets.Cine = 0x25c
		engine.MonsterOffsets.AudibleList = 0x1e8
		engine.MonsterOffsets.WaitFinished = 0x124
		engine.CineOffsets.Radius = 0x274
		engine.CineOffsets.Interruptible = 0x290
	case registry.VersionHL6153:
		engine.MonsterOffsets.MonsterState = 0x16c
		engine.MonsterOffsets.Schedule = 0x178
		engine.MonsterOffsets.ScheduleIndex = 0x17C
		engine.MonsterOffsets.Cine = 0x290
		engine.MonsterOffsets.AudibleList = 0x21C
		engine.MonsterOffsets.WaitFinished = 0x158
		engine.CineOffsets.Radius = 0x2A8
		engine.CineOffsets.Interruptible = 0x2C4
	case registry.VersionOFWON:
		engine.MonsterOffsets.MonsterState = 0x13c
		engine.MonsterOffsets.Schedule = 0x148
		engine.MonsterOffsets.ScheduleIndex = 0x14c
		engine.MonsterOffsets.Cine = 0x260
		engine.MonsterOffsets.AudibleList = 0x1ec
		engine.MonsterOffsets.WaitFinished = 0x128
		engine.CineOffsets.Radius = 0x29c
		engine.CineOffsets.Interruptible = 0x2b8
	case registry.VersionOF6153:
		engine.MonsterOffsets.MonsterState = 0x178
		engine.MonsterOffsets.Schedule = 0x184
		engine.MonsterOffsets.ScheduleIndex = 0x188
		engine.MonsterOffsets.Cine = 0x29c
		engine.MonsterOffsets.AudibleList = 0x228
		engine.MonsterOffsets.WaitFinished = 0x164
		engine.CineOffsets.Radius = 0x2D8
		engine.CineOffsets.Interruptible = 0x2F4
	case registry.VersionCSCZDS:
		engine.MonsterOffsets.MonsterState = 0x168
		engine.MonsterOffsets.Schedule = 0x174
		engine.MonsterOffsets.ScheduleIndex = 0x178
		engine.MonsterOffsets.Cine = 0x2b0
		engine.MonsterOffsets.AudibleList = 0x23c
		engine.MonsterOffsets.WaitFinished = 0x150
		engine.CineOffsets.Radius = 0x34C
		engine.CineOffsets.Interruptible = 0x368
		engine.GraphConsts.CNodeSize = 0x60
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

	switch reg.CBasePlayerForceClientDllUpdate.PatternKey() {
	case registry.VersionHL8684:
		engine.CineOffsets.Radius = 0x2AC
		engine.CineOffsets.Interruptible = 0x2C8
	case registry.VersionOF8684:
		engine.CineOffsets.Radius = 0x2DC
		engine.CineOffsets.Interruptible = 0x2F8
	case registry.VersionCSCZDS8684:
		engine.CineOffsets.Radius = 0x350
		engine.CineOffsets.Interruptible = 0x36C
	case registry.VersionCryOfFear11:
		engine.MonsterOffsets.MonsterState = 0x214
		engine.MonsterOffsets.Schedule = 0x220
		engine.MonsterOffsets.ScheduleIndex = 0x224
		engine.MonsterOffsets.Cine = 0x33C
		engine.MonsterOffsets.AudibleList = 0x2C4
		engine.MonsterOffsets.WaitFinished = 0x1FC
		engine.CineOffsets.Radius = 0x384
		engine.CineOffsets.Interruptible = 0x3AC
	case registry.VersionCryOfFear14:
		engine.MonsterOffsets.MonsterState = 0x220
		engine.MonsterOffsets.Schedule = 0x22C
		engine.MonsterOffsets.ScheduleIndex = 0x230
		engine.MonsterOffsets.Cine = 0x348
		engine.MonsterOffsets.AudibleList = 0x2D0 
		engine.MonsterOffsets.WaitFinished = 0x208
		engine.CineOffsets.Radius = 0x390
		engine.CineOffsets.Interruptible = 0x3B8
	}

	if reg.WorldGraph.Ptr() != nil {
		engine.WorldGraph.SetPointer(reg.WorldGraph.Ptr())
	}

	switch reg.CBaseMonsterRouteNew.PatternKey() {
	case registry.VersionHL8684, registry.VersionCSCZDS, registry.VersionGunman:
		engine.MonsterOffsets.Route = *(*uintptr)(unsafe.Add(reg.CBaseMonsterRouteNew.Ptr(), 0x4)) - 0xc
		engine.MonsterOffsets.RouteIndex = *(*uintptr)(unsafe.Add(reg.CBaseMonsterRouteNew.Ptr(), 0xa))
	case registry.VersionCryOfFear11:
		engine.MonsterOffsets.Route = *(*uintptr)(unsafe.Add(reg.CBaseMonsterRouteNew.Ptr(), 0x11)) - 0xc
		engine.MonsterOffsets.RouteIndex = *(*uintptr)(unsafe.Add(reg.CBaseMonsterRouteNew.Ptr(), 0x1e))
	}

	return nil
}
