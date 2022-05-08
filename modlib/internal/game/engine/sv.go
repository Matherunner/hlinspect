package engine

import "unsafe"

// SV represents the type of server_t
type SV struct {
	ptr unsafe.Pointer
}

// EntOffset returns the address offset of edict from sv.edicts
func (sv SV) EntOffset(edict uintptr) uintptr {
	edicts := *(*uintptr)(unsafe.Add(sv.ptr, 0x3bc58))
	return edict - edicts
}

// NumEdicts returns the number of edicts, sv.num_edicts
func (sv SV) NumEdicts() int {
	return *(*int)(unsafe.Add(sv.ptr, 0x3bc50))
}

// Edict returns sv.edicts[index]
func (sv SV) Edict(index int) Edict {
	base := *(*unsafe.Pointer)(unsafe.Add(sv.ptr, 0x3bc58))
	// 804 is sizeof(edict_t)
	return MakeEdict(unsafe.Add(base, index*804))
}
