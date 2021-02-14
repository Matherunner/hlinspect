package hw

type rawCVar struct {
	Name   uintptr
	String uintptr
	Flags  int32
	Value  float32
	Next   uintptr
}
