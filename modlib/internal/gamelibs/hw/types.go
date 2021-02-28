package hw

type rawCVar struct {
	Name   uintptr
	String uintptr
	Flags  int32
	Value  float32
	Next   uintptr
}

// ScreenInfo represents SCREENINFO
type ScreenInfo struct {
	Size       int32
	Width      int32
	Height     int32
	Flags      int32
	CharHeight int32
	CharWidths [256]int16
}
