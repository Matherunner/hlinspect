package engine

import "unsafe"

const (
	EdictSize = 804
)

// Edict represents edict_t
type Edict struct {
	ptr unsafe.Pointer
}

// MakeEdict creates a new instance of Edict
func MakeEdict(pointer unsafe.Pointer) Edict {
	return Edict{ptr: pointer}
}

// Ptr returns the pointer to the underlying edict_t
func (edict Edict) Ptr() unsafe.Pointer {
	return edict.ptr
}

// Free returns edict_t::free != 0
func (edict Edict) Free() bool {
	return *(*int)(edict.ptr) != 0
}

// SerialNumber returns edict_t::serialnumber
func (edict Edict) SerialNumber() int {
	return *(*int)(unsafe.Add(edict.ptr, 0x4))
}

// EntVars returns edict_t::v
func (edict Edict) EntVars() EntVars {
	return MakeEntVars(unsafe.Add(edict.ptr, 0x80))
}

// PrivateData returns edict_t::pvPrivateData
func (edict Edict) PrivateData() unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(edict.ptr, 0x7c))
}
