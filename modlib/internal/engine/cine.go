package engine

import "unsafe"

// CineOffsets contains the offsets for CCineMonster
var CineOffsets cineOffsets = cineOffsets{
	Radius:        0x2ac,
	Interruptible: 0x2c8,
}

type cineOffsets struct {
	// Search for the string "m_flRadius"
	Radius        uintptr
	Interruptible uintptr
}

// Cine represents CCineMonster
type Cine struct {
	ptr unsafe.Pointer
}

func MakeCine(pointer unsafe.Pointer) Cine {
	return Cine{ptr: pointer}
}

func (cine Cine) Pointer() unsafe.Pointer {
	return cine.ptr
}

// Radius returns CCineMonster::m_flRadius
func (cine Cine) Radius() float32 {
	return *(*float32)(unsafe.Pointer(uintptr(cine.ptr) + CineOffsets.Radius))
}

// Interruptible returns CCineMonster::m_interruptable
func (cine Cine) Interruptible() bool {
	return *(*int32)(unsafe.Pointer(uintptr(cine.ptr) + CineOffsets.Interruptible)) != 0
}
