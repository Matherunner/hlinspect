package engine

import "C"
import "unsafe"

// RawCVar represents the memory layout of cvar_t
type RawCVar struct {
	Name   uintptr
	String uintptr
	Flags  int32
	Value  float32
	Next   uintptr
}

// CVar represents cvar_t
type CVar struct {
	ptr unsafe.Pointer
}

// MakeCVar creates an instance of CVar
func MakeCVar(ptr unsafe.Pointer) CVar {
	return CVar{ptr: ptr}
}

// Pointer returns the pointer of the underlying cvar_t data
func (cvar CVar) Pointer() unsafe.Pointer {
	return cvar.ptr
}

// Name returns cvar_t::name
func (cvar CVar) Name() string {
	cstr := *(**C.char)(cvar.ptr)
	return C.GoString(cstr)
}

// Float32 returns cvar_t::value
func (cvar CVar) Float32() float32 {
	return *(*float32)(unsafe.Pointer(uintptr(cvar.ptr) + 0xc))
}

// String returns cvar_t::string
func (cvar CVar) String() string {
	cstr := *(**C.char)(unsafe.Pointer(uintptr(cvar.ptr) + 0x4))
	return C.GoString(cstr)
}
