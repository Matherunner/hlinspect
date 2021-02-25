package engine

import "unsafe"
import "C"

// Engine holds the state of the game engine
var Engine State

type GlobalVariables struct {
	address uintptr
}

func (globals *GlobalVariables) StringBase() uintptr {
	return *(*uintptr)(unsafe.Pointer(globals.address + 0x98))
}

func (globals *GlobalVariables) String(offset uint32) string {
	return C.GoString((*C.char)(unsafe.Pointer(globals.StringBase() + uintptr(offset))))
}

type EntVars struct {
	address uintptr
}

func (entvars *EntVars) Origin() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(entvars.address + 0x8))
}

func (entvars *EntVars) Classname() uint32 {
	return *(*uint32)(unsafe.Pointer(entvars.address))
}

type Edict struct {
	address uintptr
}

// Free returns free != 0
func (edict *Edict) Free() bool {
	return *(*int)(unsafe.Pointer(edict.address)) != 0
}

func (edict *Edict) EntVars() *EntVars {
	address := uintptr(unsafe.Pointer(edict.address + 0x80))
	return &EntVars{address: address}
}

// SV represents the type of server_t
type SV struct {
	address uintptr
}

// NumEdicts returns the number of edicts, sv.num_edicts
func (sv *SV) NumEdicts() int {
	return *(*int)(unsafe.Pointer(sv.address + 0x3bc58))
}

// Edict returns the address of the edict_t, or sv.edict[index]
func (sv *SV) Edict(index int) *Edict {
	base := *(*uintptr)(unsafe.Pointer(sv.address + 0x3bc60))
	// 804 is sizeof(edict_t)
	return &Edict{address: base + uintptr(index*804)}
}

// State interface to the game engine
type State struct {
	SV              SV
	GlobalVariables GlobalVariables
	ppmove          uintptr
}

func (eng *State) SetGlobalVariables(address uintptr) {
	eng.GlobalVariables.address = address
}

// SetSV sets the address of sv
func (eng *State) SetSV(address uintptr) {
	eng.SV.address = address
}

// SetPPMove sets the address of ppmove
func (eng *State) SetPPMove(address uintptr) {
	eng.ppmove = address
}

// PMoveVelocity returns pmove->velocity
func (eng *State) PMoveVelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0x5c))
}

// PMovePosition returns pmove->origin
func (eng *State) PMovePosition() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0x38))
}

// PMoveViewangles returns pmove->angles
func (eng *State) PMoveViewangles() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0x44))
}

// PMoveBasevelocity returns pmove->basevelocity
func (eng *State) PMoveBasevelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0x74))
}

// PMoveCmdFSU returns pmove->cmd.forwardmove treated as a 3-element float array
func (eng *State) PMoveCmdFSU() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0x45468))
}

// PMovePunchangles returns pmove->punchangle
func (eng *State) PMovePunchangles() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(eng.ppmove + 0xa0))
}

// PMoveEntFriction returns pmove->friction
func (eng *State) PMoveEntFriction() float32 {
	return *(*float32)(unsafe.Pointer(eng.ppmove + 0xc4))
}

// PMoveEntGravity returns pmove->gravity
func (eng *State) PMoveEntGravity() float32 {
	return *(*float32)(unsafe.Pointer(eng.ppmove + 0xc0))
}

// PMoveFrameTime returns pmove->cmd.msec
func (eng *State) PMoveFrameTime() uint32 {
	return uint32(*(*uint8)(unsafe.Pointer(eng.ppmove + 0x4545a)))
}

// PMoveCmdButtons returns pmove->cmd.buttons
func (eng *State) PMoveCmdButtons() uint32 {
	return uint32(*(*uint16)(unsafe.Pointer(eng.ppmove + 0x45476)))
}

// PMoveOnground returns the boolean value of (pmove->onground != -1)
func (eng *State) PMoveOnground() bool {
	return *(*int32)(unsafe.Pointer(eng.ppmove + 0xe0)) != -1
}

// PMoveFlags returns pmove->flags
func (eng *State) PMoveFlags() uint32 {
	return *(*uint32)(unsafe.Pointer(eng.ppmove + 0xb8))
}

// PMoveWaterlevel returns pmove->waterlevel
func (eng *State) PMoveWaterlevel() uint32 {
	return *(*uint32)(unsafe.Pointer(eng.ppmove + 0xe4))
}

// PMoveInDuck returns (pmove->bInDuck != 0)
func (eng *State) PMoveInDuck() bool {
	return *(*int32)(unsafe.Pointer(eng.ppmove + 0x90)) != 0
}

func (eng *State) PMoveLadder() bool {
	// TODO: get global variable
	return false
}

// PMoveImpulse returns pmove->cmd.impulse
func (eng *State) PMoveImpulse() uint32 {
	return uint32(*(*uint8)(unsafe.Pointer(eng.ppmove + 0x45478)))
}
