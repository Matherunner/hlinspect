package engine

import (
	"unsafe"
)

// MonsterOffsets store offsets to class members
var MonsterOffsets monsterOffsets = monsterOffsets{
	Schedule:      0x178,
	ScheduleIndex: 0x17c,
	Cine:          0x290,
}

type monsterOffsets struct {
	// Look inside CBaseMonster::ChangeSchedule.
	Schedule      uintptr
	ScheduleIndex uintptr
	// Look inside CBaseMonster::GetScheduleOfType. Search for "Script failed for %s"
	Cine uintptr
}

// Monster represents CBaseMonster
type Monster struct {
	address uintptr
}

// MakeMonster creates a new instance of Monster
func MakeMonster(address uintptr) Monster {
	return Monster{address: address}
}

// Schedule returns CBaseMonster::m_pSchedule
func (monster Monster) Schedule() *Schedule {
	address := *(*uintptr)(unsafe.Pointer(monster.address + MonsterOffsets.Schedule))
	if address == 0 {
		return nil
	}
	schedule := MakeSchedule(address)
	return &schedule
}

// ScheduleIndex returns CBaseMonster::m_iScheduleIndex
func (monster Monster) ScheduleIndex() int {
	return int(*(*int32)(unsafe.Pointer(monster.address + MonsterOffsets.ScheduleIndex)))
}

// Cine returns CBaseMonster::m_pCine
func (monster Monster) CineAddr() uintptr {
	return *(*uintptr)(unsafe.Pointer(monster.address + MonsterOffsets.Cine))
}
