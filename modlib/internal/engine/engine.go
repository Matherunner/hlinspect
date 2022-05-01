package engine

import "unsafe"
import "C"

// Engine holds the state of the game engine
var Engine State

type GlobalVariables struct {
	ptr unsafe.Pointer
}

// Time returns globalvars_t::time
func (globals GlobalVariables) Time() float32 {
	return *(*float32)(globals.ptr)
}

func (globals GlobalVariables) StringBase() unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(uintptr(globals.ptr) + 0x98))
}

func (globals GlobalVariables) String(offset uint32) string {
	return C.GoString((*C.char)(unsafe.Pointer(uintptr(globals.StringBase()) + uintptr(offset))))
}

// EntVars represents entvars_t
type EntVars struct {
	ptr unsafe.Pointer
}

// MakeEntVars creates an instance of EntVars
func MakeEntVars(pointer unsafe.Pointer) EntVars {
	return EntVars{ptr: pointer}
}

// Pointer returns the pointer to this entvars_t object
func (entvars EntVars) Pointer() unsafe.Pointer {
	return entvars.ptr
}

// Origin returns entvars_t::origin
func (entvars EntVars) Origin() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(entvars.ptr) + 0x8))
}

// Angles returns entvars_t::angles
func (entvars EntVars) Angles() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(entvars.ptr) + 0x50))
}

// Classname returns entvars_t::classname
func (entvars EntVars) Classname() uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(entvars.ptr) + 0x0))
}

// Targetname returns entvars_t::targetname
func (entvars EntVars) Targetname() uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(entvars.ptr) + 0x1cc))
}

// Mins returns entvars_t::mins
func (entvars EntVars) Mins() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(entvars.ptr) + 0xdc))
}

// Maxs returns entvars_t::maxs
func (entvars EntVars) Maxs() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(entvars.ptr) + 0xe8))
}

// AbsMin returns entvars_t::absmin
func (entvars EntVars) AbsMin() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(entvars.ptr) + 0xc4))
}

// AbsMax returns entvars_t::absmax
func (entvars EntVars) AbsMax() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(entvars.ptr) + 0xd0))
}

// Edict represents edict_t
type Edict struct {
	ptr unsafe.Pointer
}

// MakeEdict creates a new instance of Edict
func MakeEdict(pointer unsafe.Pointer) Edict {
	return Edict{ptr: pointer}
}

// Pointer returns the pointer to the underlying edict_t
func (edict Edict) Pointer() unsafe.Pointer {
	return edict.ptr
}

// Free returns free != 0
func (edict Edict) Free() bool {
	return *(*int)(edict.ptr) != 0
}

// SerialNumber returns edict_t::serialnumber
func (edict Edict) SerialNumber() int {
	return *(*int)(unsafe.Pointer(uintptr(edict.ptr) + 0x4))
}

// EntVars returns edict_t::v
func (edict Edict) EntVars() EntVars {
	return MakeEntVars(unsafe.Pointer(uintptr(edict.ptr) + 0x80))
}

// PrivateData returns edict_t::pvPrivateData
func (edict Edict) PrivateData() unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(uintptr(edict.ptr) + 0x7c))
}

// SV represents the type of server_t
type SV struct {
	ptr unsafe.Pointer
}

// EntOffset returns the address offset of edict from sv.edicts
func (sv SV) EntOffset(edict uintptr) uintptr {
	edicts := *(*uintptr)(unsafe.Pointer(uintptr(sv.ptr) + 0x3bc60))
	return edict - edicts
}

// NumEdicts returns the number of edicts, sv.num_edicts
func (sv SV) NumEdicts() int {
	return *(*int)(unsafe.Pointer(uintptr(sv.ptr) + 0x3bc58))
}

// Edict returns sv.edicts[index]
func (sv SV) Edict(index int) Edict {
	base := *(*unsafe.Pointer)(unsafe.Add(sv.ptr, 0x3bc60))
	// 804 is sizeof(edict_t)
	return MakeEdict(unsafe.Add(base, index*804))
}

// State interface to the game engine
type State struct {
	SV              SV
	GlobalVariables GlobalVariables
	ppmove          unsafe.Pointer
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
	return *(*[3]float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0x5c))
}

// PMovePosition returns pmove->origin
func (eng State) PMovePosition() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0x38))
}

// PMoveViewangles returns pmove->angles
func (eng State) PMoveViewangles() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0x44))
}

// PMoveBasevelocity returns pmove->basevelocity
func (eng State) PMoveBasevelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0x74))
}

// PMoveCmdFSU returns pmove->cmd.forwardmove treated as a 3-element float array
func (eng State) PMoveCmdFSU() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0x45468))
}

// PMovePunchangles returns pmove->punchangle
func (eng State) PMovePunchangles() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0xa0))
}

// PMoveEntFriction returns pmove->friction
func (eng State) PMoveEntFriction() float32 {
	return *(*float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0xc4))
}

// PMoveEntGravity returns pmove->gravity
func (eng State) PMoveEntGravity() float32 {
	return *(*float32)(unsafe.Pointer(uintptr(eng.ppmove) + 0xc0))
}

// PMoveFrameTime returns pmove->cmd.msec
func (eng State) PMoveFrameTime() uint32 {
	return uint32(*(*uint8)(unsafe.Pointer(uintptr(eng.ppmove) + 0x4545a)))
}

// PMoveCmdButtons returns pmove->cmd.buttons
func (eng State) PMoveCmdButtons() uint32 {
	return uint32(*(*uint16)(unsafe.Pointer(uintptr(eng.ppmove) + 0x45476)))
}

// PMoveOnground returns the boolean value of (pmove->onground != -1)
func (eng State) PMoveOnground() bool {
	return *(*int32)(unsafe.Pointer(uintptr(eng.ppmove) + 0xe0)) != -1
}

// PMoveFlags returns pmove->flags
func (eng State) PMoveFlags() uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(eng.ppmove) + 0xb8))
}

// PMoveWaterlevel returns pmove->waterlevel
func (eng State) PMoveWaterlevel() uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(eng.ppmove) + 0xe4))
}

// PMoveInDuck returns (pmove->bInDuck != 0)
func (eng State) PMoveInDuck() bool {
	return *(*int32)(unsafe.Pointer(uintptr(eng.ppmove) + 0x90)) != 0
}

func (eng State) PMoveLadder() bool {
	// TODO: get global variable
	return false
}

// PMoveImpulse returns pmove->cmd.impulse
func (eng State) PMoveImpulse() uint32 {
	return uint32(*(*uint8)(unsafe.Pointer(uintptr(eng.ppmove) + 0x45478)))
}
