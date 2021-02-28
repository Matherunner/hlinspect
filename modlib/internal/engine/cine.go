package engine

import "unsafe"

var CineOffsets cineOffsets = cineOffsets{
	Radius: 0x2ac,
}

type cineOffsets struct {
	Radius uintptr
}

type Cine struct {
	address uintptr
}

func MakeCine(address uintptr) Cine {
	return Cine{address: address}
}

func (cine Cine) Radius() float32 {
	return *(*float32)(unsafe.Pointer(cine.address + CineOffsets.Radius))
}
