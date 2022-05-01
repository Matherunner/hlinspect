package gamelibs

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs/cdefs"
	"hlinspect/internal/hooks"
	"unsafe"
)

var hlDLL *hooks.Module

func initHLRegistry(reg *APIRegistry) {
	reg.PMInit = hooks.MakeFunctionPattern("PM_Init", map[string]string{"Windows": "PM_Init"}, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 55 08 33 C0 56 8D 8A ?? ?? ?? ?? 8B B0 ?? ?? ?? ?? 83 C0 0C 89 71 FC 8B B0 ?? ?? ?? ?? 89 31"),
		OF8684: hooks.MustMakePattern("8B 44 24 04 A3 ?? ?? ?? ?? E8 ?? ?? ?? ?? E8 ?? ?? ?? ?? C7 05"),
	})
	reg.PMPlayerMove = hooks.MakeFunctionPattern("PM_PlayerMove", map[string]string{"Windows": "PM_PlayerMove"}, map[string]hooks.SearchPattern{
		HL8684:     hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 57 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 BA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
		BigLolly:   hooks.MustMakePattern("55 8B EC 83 EC 0C C7 45 FC 00 00 00 00 A1 ?? ?? ?? ?? 8B 4D 08 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? C7 82 8C 54 04 00 00 00 00 00 A1"),
		TWHLTower2: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 8B 4D 08 53 56 57 33 FF 89 7D FC 89 48 04 E8 D8 FC FF FF A1 ?? ?? ?? ?? 89 B8 8C 54 04 00 A1 ?? ?? ?? ?? 0F B6 88 5A 54 04 00"),
		CSCZDS:     hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 56 57 33 ED 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 AA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
	})
	reg.CSoundEntActiveList = hooks.MakeFunctionPattern("CSoundEnt::ActiveList", map[string]string{WindowsHLDLL: "CSoundEnt::ActiveList"}, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 58 C3"),
		OF8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 64 C3"),
		HLWON:  hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 24 C3"),
		CSCZDS: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 50 C3"),
	})
	reg.CSoundEntSoundPointerForIndex = hooks.MakeFunctionPattern("CSoundEnt::SoundPointerForIndex", map[string]string{WindowsHLDLL: "CSoundEnt::SoundPointerForIndex"}, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("8B 0D ?? ?? ?? ?? 85 C9 75 03 33 C0 C3 8B 44 24 04 83 F8 3F 7E 13 68 ?? ?? ?? ?? 6A 01 FF 15 ?? ?? ?? ??"),
	})
	reg.CBaseMonsterChangeSchedule = hooks.MakeFunctionPattern("CBaseMonster::ChangeSchedule", map[string]string{WindowsHLDLL: "CBaseMonster::ChangeSchedule"}, map[string]hooks.SearchPattern{
		// Search for COND_HEAR_SOUND
		HL8684: hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 78 01 00 00 89 91 7C 01 00 00 89 91 74 01 00 00 89 91 F0 00 00 00 89 91 68 02 00 00"),
		OF8684: hooks.MustMakePattern("8B 81 84 01 00 00 33 D2 3B C2 56 74 55 8B 00 3B C2 74 4F 8B B1 88 01 00 00 57 8B 3C F0"),
		HLWON:  hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 44 01 00 00 89 91 48 01 00 00 89 91 40 01 00 00 89 91 BC 00 00 00 89 91 34 02 00 00"),
		OFWON:  hooks.MustMakePattern("8B 81 48 01 00 00 33 D2 3B C2 56 74 55 8B 00 3B C2 74 4F 8B B1 4C 01 00 00 57 8B 3C F0"),
		CSCZDS: hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 74 01 00 00 89 91 78 01 00 00 89 91 70 01 00 00 89 91 7C 01 00 00 89 91 88 02 00 00"),
		Gunman: hooks.MustMakePattern("8B 44 24 04 53 57 8B F9 33 DB 89 87 4C 01 00 00 89 9F 50 01 00 00 89 9F 48 01 00 00 89 9F BC 00 00 00 89 9F 3C 02 00 00"),
	})
	reg.CBaseMonsterRouteNew = hooks.MakeFunctionPattern("CBaseMonster::RouteNew", map[string]string{
		WindowsHLDLL: "CBaseMonster::RouteNew",
	}, map[string]hooks.SearchPattern{
		// Search for "No Path from %d to %d!" to find CBaseMonster::FGetNodeRoute
		// Go to any of the cross reference, the first function should be CBaseMonster::RouteNew
		// This pattern includes the initial part of CBaseMonster::FRouteClear, and mask out the offsets
		HL8684: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? C1 E0 04"),
		CSCZDS: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? 83 C0 14"),
		Gunman: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? 83 C0 16"),
	})
	reg.CBaseMonsterPBestSound = hooks.MakeFunctionPattern("CBaseMonster::PBestSound", map[string]string{
		WindowsHLDLL: "CBaseMonster::PBestSound",
	}, map[string]hooks.SearchPattern{
		// Search for "ERROR! monster %s has no audible sounds!"
		HL8684: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB 1C 02 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
		OF8684: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB 28 02 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
		HLWON:  hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB E8 01 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
		OFWON:  hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB EC 01 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
		CSCZDS: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB 3C 02 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
		Gunman: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB F0 01 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
	})
	reg.WorldGraph = hooks.MakeFunctionPattern("WorldGraph", map[string]string{
		// Not actually a function
		WindowsHLDLL: "WorldGraph",
	}, nil)
	reg.CGraphInitGraph = hooks.MakeFunctionPattern("CGraph::InitGraph", map[string]string{
		WindowsHLDLL: "CGraph::InitGraph",
	}, map[string]hooks.SearchPattern{
		// Search for "Couldn't malloc %d nodes!" to find CGraph::AllocNodes
		// Then find cross reference from CWorld::Precache
		HL8684: hooks.MustMakePattern("56 8B F1 57 33 FF 8B 46 10 89 3E 3B C7 89 7E 04 89 7E 08 74 0C 50 E8 ?? ?? ?? ?? 83 C4 04 89 7E 10 8B 46 0C"),
	})
}

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

	initHLRegistry(reg)

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
	case OF8684:
		engine.MonsterOffsets.MonsterState = 0x178
		engine.MonsterOffsets.Schedule = 0x184
		engine.MonsterOffsets.ScheduleIndex = 0x188
		engine.MonsterOffsets.Cine = 0x29c
		engine.MonsterOffsets.AudibleList = 0x228
		engine.MonsterOffsets.WaitFinished = 0x164
		engine.CineOffsets.Radius = 0x2dc
		engine.CineOffsets.Interruptible = 0x2f8
	case HLWON:
		engine.MonsterOffsets.MonsterState = 0x138
		engine.MonsterOffsets.Schedule = 0x144
		engine.MonsterOffsets.ScheduleIndex = 0x148
		engine.MonsterOffsets.Cine = 0x25c
		engine.MonsterOffsets.AudibleList = 0x1e8
		engine.MonsterOffsets.WaitFinished = 0x124
		engine.CineOffsets.Radius = 0x274
		engine.CineOffsets.Interruptible = 0x290
	case OFWON:
		engine.MonsterOffsets.MonsterState = 0x13c
		engine.MonsterOffsets.Schedule = 0x148
		engine.MonsterOffsets.ScheduleIndex = 0x14c
		engine.MonsterOffsets.Cine = 0x260
		engine.MonsterOffsets.AudibleList = 0x1ec
		engine.MonsterOffsets.WaitFinished = 0x128
		engine.CineOffsets.Radius = 0x29c
		engine.CineOffsets.Interruptible = 0x2b8
	case CSCZDS:
		engine.MonsterOffsets.MonsterState = 0x168
		engine.MonsterOffsets.Schedule = 0x174
		engine.MonsterOffsets.ScheduleIndex = 0x178
		engine.MonsterOffsets.Cine = 0x2b0
		engine.MonsterOffsets.AudibleList = 0x23c
		engine.MonsterOffsets.WaitFinished = 0x150
		engine.CineOffsets.Radius = 0x350
		engine.CineOffsets.Interruptible = 0x36c
	case Gunman:
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
	case WindowsHLDLL:
		engine.CineOffsets.Radius = 0x2a8
		engine.CineOffsets.Interruptible = 0x2c4
	}

	if reg.WorldGraph.Address() != nil {
		engine.WorldGraph.SetPointer(reg.WorldGraph.Address())
	}

	switch reg.CBaseMonsterRouteNew.PatternKey() {
	case HL8684, CSCZDS, Gunman:
		engine.MonsterOffsets.Route = *(*uintptr)(unsafe.Pointer(uintptr(reg.CBaseMonsterRouteNew.Address()) + 0x4)) - 0xc
		engine.MonsterOffsets.RouteIndex = *(*uintptr)(unsafe.Pointer(uintptr(reg.CBaseMonsterRouteNew.Address()) + 0xa))
	}

	return nil
}
