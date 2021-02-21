package engine

import "unsafe"

// Engine holds the state of the game engine
var Engine State

// State interface to the game engine
type State struct {
	ppmove uintptr
}

// SetPPMove set ppmove address
func (eng *State) SetPPMove(ppmove uintptr) {
	eng.ppmove = ppmove
}

func (eng *State) PMoveVelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0x5c))
}
