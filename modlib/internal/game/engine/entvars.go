package engine

import "unsafe"

// EntVars represents entvars_t
type EntVars struct {
	ptr unsafe.Pointer
}

// MakeEntVars creates an instance of EntVars
func MakeEntVars(pointer unsafe.Pointer) EntVars {
	return EntVars{ptr: pointer}
}

// Ptr returns the pointer to this entvars_t object
func (entvars EntVars) Ptr() unsafe.Pointer {
	return entvars.ptr
}

// Origin returns entvars_t::origin
func (entvars EntVars) Origin() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0x8))
}

// Angles returns entvars_t::angles
func (entvars EntVars) Angles() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0x50))
}

// Classname returns entvars_t::classname
func (entvars EntVars) Classname() uint32 {
	return *(*uint32)(unsafe.Add(entvars.ptr, 0x0))
}

// Targetname returns entvars_t::targetname
func (entvars EntVars) Targetname() uint32 {
	return *(*uint32)(unsafe.Add(entvars.ptr, 0x1cc))
}

// Mins returns entvars_t::mins
func (entvars EntVars) Mins() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xdc))
}

// Maxs returns entvars_t::maxs
func (entvars EntVars) Maxs() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xe8))
}

// AbsMin returns entvars_t::absmin
func (entvars EntVars) AbsMin() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xc4))
}

// AbsMax returns entvars_t::absmax
func (entvars EntVars) AbsMax() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xd0))
}
