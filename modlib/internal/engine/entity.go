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
	ptr unsafe.Pointer
}

// MakeEntity creates an instance of Entity
func MakeEntity(pointer unsafe.Pointer) Entity {
	return Entity{ptr: pointer}
}

// EntVars returns CBaseEntity::pev
func (entity Entity) EntVars() EntVars {
	return MakeEntVars(unsafe.Pointer(uintptr(entity.ptr) + EntityOffsets.PEV))
}
