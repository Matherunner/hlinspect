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
	address uintptr
}

func MakeCine(address uintptr) Cine {
	return Cine{address: address}
}

func (cine Cine) Radius() float32 {
	return *(*float32)(unsafe.Pointer(cine.address + CineOffsets.Radius))
}

// Interruptible returns CCineMonster::m_interruptable
func (cine Cine) Interruptible() bool {
	return *(*int32)(unsafe.Pointer(cine.address + CineOffsets.Interruptible)) != 0
}
