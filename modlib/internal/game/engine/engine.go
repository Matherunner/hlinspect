package engine

import "unsafe"
import "C"

type GlobalVariables struct {
	ptr unsafe.Pointer
}

// Time returns globalvars_t::time
func (globals GlobalVariables) Time() float32 {
	return *(*float32)(globals.ptr)
}

func (globals GlobalVariables) StringBase() unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(globals.ptr, 0x98))
}

func (globals GlobalVariables) String(offset uint32) string {
	return C.GoString((*C.char)(unsafe.Add(globals.StringBase(), offset)))
}

// State interface to the game engine
type State struct {
	SV              SV
	GlobalVariables GlobalVariables
	ppmove          unsafe.Pointer
}

func NewState() *State {
	return &State{}
}

func (eng *State) SetGlobalVariables(pointer unsafe.Pointer) {
	eng.GlobalVariables.ptr = pointer
}

// SetSV sets the address of sv
func (eng *State) SetSV(pointer unsafe.Pointer) {
	eng.SV.ptr = pointer
}

// SetPPMove sets the address of ppmove
func (eng *State) SetPPMove(pointer unsafe.Pointer) {
	eng.ppmove = pointer
}

// PMoveVelocity returns pmove->velocity
func (eng State) PMoveVelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Add(eng.ppmove, 0x5c))
}

// PMovePosition returns pmove->origin
func (eng State) PMovePosition() [3]float32 {
	return *(*[3]float32)(unsafe.Add(eng.ppmove, 0x38))
}

// PMoveViewangles returns pmove->angles
func (eng State) PMoveViewangles() [3]float32 {
	return *(*[3]float32)(unsafe.Add(eng.ppmove, 0x44))
}

// PMoveBasevelocity returns pmove->basevelocity
func (eng State) PMoveBasevelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Add(eng.ppmove, 0x74))
}

// PMoveCmdFSU returns pmove->cmd.forwardmove treated as a 3-element float array
func (eng State) PMoveCmdFSU() [3]float32 {
	return *(*[3]float32)(unsafe.Add(eng.ppmove, 0x45468))
}

// PMovePunchangles returns pmove->punchangle
func (eng State) PMovePunchangles() [3]float32 {
	return *(*[3]float32)(unsafe.Add(eng.ppmove, 0xa0))
}

// PMoveEntFriction returns pmove->friction
func (eng State) PMoveEntFriction() float32 {
	return *(*float32)(unsafe.Add(eng.ppmove, 0xc4))
}

// PMoveEntGravity returns pmove->gravity
func (eng State) PMoveEntGravity() float32 {
	return *(*float32)(unsafe.Add(eng.ppmove, 0xc0))
}

// PMoveFrameTime returns pmove->cmd.msec
func (eng State) PMoveFrameTime() uint32 {
	return uint32(*(*uint8)(unsafe.Add(eng.ppmove, 0x4545a)))
}

// PMoveCmdButtons returns pmove->cmd.buttons
func (eng State) PMoveCmdButtons() uint32 {
	return uint32(*(*uint16)(unsafe.Add(eng.ppmove, 0x45476)))
}

// PMoveOnground returns the boolean value of (pmove->onground != -1)
func (eng State) PMoveOnground() bool {
	return *(*int32)(unsafe.Add(eng.ppmove, 0xe0)) != -1
}

// PMoveFlags returns pmove->flags
func (eng State) PMoveFlags() uint32 {
	return *(*uint32)(unsafe.Add(eng.ppmove, 0xb8))
}

// PMoveWaterlevel returns pmove->waterlevel
func (eng State) PMoveWaterlevel() uint32 {
	return *(*uint32)(unsafe.Add(eng.ppmove, 0xe4))
}

// PMoveInDuck returns (pmove->bInDuck != 0)
func (eng State) PMoveInDuck() bool {
	return *(*int32)(unsafe.Add(eng.ppmove, 0x90)) != 0
}

func (eng State) PMoveLadder() bool {
	// TODO: get global variable
	return false
}

// PMoveImpulse returns pmove->cmd.impulse
func (eng State) PMoveImpulse() uint32 {
	return uint32(*(*uint8)(unsafe.Add(eng.ppmove, 0x45478)))
}
