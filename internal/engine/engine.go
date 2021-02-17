package engine

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
