package engine

import (
	"unsafe"
)

// MonsterOffsets store offsets to class members
var MonsterOffsets monsterOffsets = monsterOffsets{
	Schedule:      0x178,
	ScheduleIndex: 0x17c,
	Cine:          0x290,
	Route:         0x180,
	RouteIndex:    0x204,
}

type monsterOffsets struct {
	// Look inside CBaseMonster::ChangeSchedule.
	Schedule      uintptr
	ScheduleIndex uintptr
	// Look inside CBaseMonster::GetScheduleOfType. Search for "Script failed for %s"
	Cine uintptr
	// Found from CBaseMonster::RouteNew
	Route      uintptr
	RouteIndex uintptr
}

// Monster represents CBaseMonster
type Monster struct {
	ptr unsafe.Pointer
}

// MakeMonster creates a new instance of Monster
func MakeMonster(pointer unsafe.Pointer) Monster {
	return Monster{ptr: pointer}
}

// Schedule returns CBaseMonster::m_pSchedule
func (monster Monster) Schedule() *Schedule {
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(monster.ptr) + MonsterOffsets.Schedule))
	if ptr == nil {
		return nil
	}
	schedule := MakeSchedule(uintptr(ptr))
	return &schedule
}

// ScheduleIndex returns CBaseMonster::m_iScheduleIndex
func (monster Monster) ScheduleIndex() int {
	return int(*(*int32)(unsafe.Pointer(uintptr(monster.ptr) + MonsterOffsets.ScheduleIndex)))
}

// Cine returns CBaseMonster::m_pCine
func (monster Monster) Cine() Cine {
	return MakeCine(*(*unsafe.Pointer)(unsafe.Pointer(uintptr(monster.ptr) + MonsterOffsets.Cine)))
}
