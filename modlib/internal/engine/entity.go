package engine

import "unsafe"

// EntityOffsets contains offsets for CBaseEntity
var EntityOffsets entityOffsets = entityOffsets{
	PEV: 0x4,
}

type entityOffsets struct {
	PEV uintptr
}

// Entity represents CBaseEntity
type Entity struct {
	address uintptr
}

// MakeEntity creates an instance of Entity
func MakeEntity(address uintptr) Entity {
	return Entity{address: address}
}

// EntVars returns CBaseEntity::pev
func (entity Entity) EntVars() EntVars {
	return MakeEntVars(unsafe.Pointer(entity.address + EntityOffsets.PEV))
}
