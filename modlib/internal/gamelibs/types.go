package gamelibs

import "unsafe"

const (
	TriTriangles = iota
	TriTriangleFan
	TriQuads
	TriPolygon
	TriLines
	TriTriangleStrip
	TriQuadStrip
)

const (
	TriFront = iota
	TriNone
)

const (
	KRenderNormal = iota
	KRenderTransColor
	KRenderTransTexture
	KRenderGlow
	KRenderTransAlpha
	KRenderTransAdd
)

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

const (
	TraceDontIgnoreMonsters = iota
	TraceIgnoreMonsters
	TraceMissile
)

// TraceResult represents TraceResult
type TraceResult struct {
	AllSolid    int32
	StartSolid  int32
	InOpen      int32
	InWater     int32
	Fraction    float32
	EndPos      [3]float32
	PlaneDist   float32
	PlaneNormal [3]float32
	Hit         unsafe.Pointer
	HitGroup    int32
}
