package hw

import (
	"hlinspect/internal/cvar"
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gl"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"unsafe"
)

/*
#include <stdlib.h>
#include "defs.h"
*/
import "C"

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

var hwDLL *hooks.Module

var buildNumberPattern = hooks.MakeFunctionPattern("build_number", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
})
var cvarRegisterVariablePattern = hooks.MakeFunctionPattern("Cvar_RegisterVariable", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 14 53 56 8B 75 08 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68"),
	gamelibs.HLNGHL: hooks.MustMakePattern("83 EC 14 53 56 8B 74 24 20 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5B 83 C4 14 C3 8B 16 52 E8"),
})
var vFadeAlphaPattern = hooks.MakeFunctionPattern("V_FadeAlpha", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? DF E0 F6 C4 05 7A 1C D9 05 ?? ?? ?? ?? DC 1D"),
	gamelibs.HL4554: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? 83 EC"),
})
var drawStringPattern = hooks.MakeFunctionPattern("Draw_String", nil, map[string]hooks.SearchPattern{
	// Search for "%i %i %i", there is a thunk that call this function. There are two functions that call the thunk.
	// The shorter one is Draw_String.
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 56 57 E8 ?? ?? ?? ?? 8B 4D 0C 8B 75 08 50 8B 45 10 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
})
var vgui2DrawSetTextColorAlphaPattern = hooks.MakeFunctionPattern("VGUI2_Draw_SetTextColorAlpha", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8A 45 08 8A 4D 0C 8A 55 10 88 45 08 8A 45 14 88 4D 09 88 55 0A 88 45 0B 8B 4D 08 89"),
})
var hostAutoSaveFPattern = hooks.MakeFunctionPattern("Host_AutoSave_f", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? B9 01 00 00 00 3B C1 0F 85 9F 00 00 00 A1 ?? ?? ?? ?? 85 C0 75 10 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 04 33 C0 C3 39 0D"),
})
var hostNoclipFPattern = hooks.MakeFunctionPattern("Host_Noclip_f", nil, map[string]hooks.SearchPattern{
	// Search for "noclip ON\n"
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 24 A1 ?? ?? ?? ?? BA 01 00 00 00 3B C2 75 09 E8 ?? ?? ?? ?? 8B E5 5D C3 D9 05 ?? ?? ?? ?? D8 1D"),
})
var pfTracelineDLLPattern = hooks.MakeFunctionPattern("PF_traceline_DLL", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 14 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4D 0C 8B 55 08 56 50 8B 45 10 50 51 52 E8 ?? ?? ?? ?? D9 05"),
})
var triGLRenderModePattern = hooks.MakeFunctionPattern("tri_GL_RenderMode", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 56 8B 75 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
})
var triGLBeginPattern = hooks.MakeFunctionPattern("tri_GL_Begin", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 45 08 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? 5D C3"),
})
var triGLEndPattern = hooks.MakeFunctionPattern("tri_GL_End", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 55 8B EC 8B 45 0C"),
})
var triGLColor4fPattern = hooks.MakeFunctionPattern("tri_GL_Color4f", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 51 83 3D ?? ?? ?? ?? 04 75 4A D9 45 14 D8 0D ?? ?? ?? ?? D9 5D FC D9 45 FC E8 ?? ?? ?? ?? D9 45 10"),
})
var triGLCullFacePattern = hooks.MakeFunctionPattern("tri_GL_CullFace", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 83 E8 00 74 10 48 75 23 68 44 0B 00 00 FF 15 ?? ?? ?? ?? 5D C3 68 44 0B 00 00"),
})
var triGLVertex3fvPattern = hooks.MakeFunctionPattern("tri_GL_Vertex3fv", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 50 FF 15 ?? ?? ?? ?? 5D C3 90 55 8B EC 8B 45 10 8B 4D 0C 8B 55 08 50 51 52"),
})
var screenTransformPattern = hooks.MakeFunctionPattern("ScreenTransform", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 51 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
})
var worldTransformPattern = hooks.MakeFunctionPattern("WorldTransform", nil, map[string]hooks.SearchPattern{
	// Most likely the function below ScreenTransform
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48"),
})
var hudGetScreenInfoPattern = hooks.MakeFunctionPattern("hudGetScreenInfo", nil, map[string]hooks.SearchPattern{
	// Search for "Half-Life %i/%s (hw build %d)". This function is Draw_ConsoleBackground
	// The function below it should be Draw_FillRGBA. Get the cross references to Draw_FillRGBA. One of them
	// is a global variable of enginefuncs. The next entry is hudGetScreenInfo.
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 83 C4 04 85 C0 75 02 5d"),
})
var rDrawSequentialPolyPattern = hooks.MakeFunctionPattern("R_DrawSequentialPoly", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 53 56 57 83 B8 F8 02 00 00 01 75 63 E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
	gamelibs.HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 53 55 56 8B 88 F8 02 00 00 BE 01 00 00 00 3B CE 57 75 62 E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
})
var rClearPattern = hooks.MakeFunctionPattern("R_Clear", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("8B 15 ?? ?? ?? ?? 33 C0 83 FA 01 0F 9F C0 50 E8 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 83 C4 04 DF E0"),
})
var memoryInitPattern = hooks.MakeFunctionPattern("Memory_Init", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 8B 4D 0C 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
})
var pfCheckClientIPattern = hooks.MakeFunctionPattern("PF_checkclient_I", nil, map[string]hooks.SearchPattern{
	// Search for "Spawned a NULL entity!", the referencing function is CreateNamedEntity
	// Find cross references, go to the global data, that data is g_engfuncsExportedToDlls
	// Go up 6 entries, we will end up at PF_checkclient_I
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 0C DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
})

// GetScreenInfo proxies hudGetScreenInfo
func GetScreenInfo() ScreenInfo {
	screenInfo := ScreenInfo{}
	screenInfo.Size = int32(unsafe.Sizeof(screenInfo))
	hooks.CallFuncInts1(hudGetScreenInfoPattern.Address(), uintptr(unsafe.Pointer(&screenInfo)))
	return screenInfo
}

// TraceLine traces a line and return the hit results
func TraceLine() {
	// TODO:
}

// BuildNumber build_number
func BuildNumber() int {
	return hooks.CallFuncInts0(buildNumberPattern.Address())
}

// CvarRegisterVariable Cvar_RegisterVariable
func CvarRegisterVariable(cvar uintptr) {
	hooks.CallFuncInts1(cvarRegisterVariablePattern.Address(), cvar)
}

// DrawString Draw_String
func DrawString(x, y int, text string) {
	ctext := unsafe.Pointer(C.CString(text))
	defer C.free(ctext)
	hooks.CallFuncInts3(drawStringPattern.Address(), uintptr(x), uintptr(y), uintptr(ctext))
}

// VGUI2DrawSetTextColorAlpha VGUI2_Draw_SetTextColorAlpha
func VGUI2DrawSetTextColorAlpha(r, g, b, a int) {
	hooks.CallFuncInts4(vgui2DrawSetTextColorAlphaPattern.Address(), uintptr(r), uintptr(g), uintptr(b), uintptr(a))
}

// HookedVFadeAlpha V_FadeAlpha
//export HookedVFadeAlpha
func HookedVFadeAlpha() int {
	if cvar.FadeRemove.Float32() != 0 {
		return 0
	}
	return hooks.CallFuncInts0(vFadeAlphaPattern.Address())
}

// HookedRClear R_Clear
//export HookedRClear
func HookedRClear() {
	if cvar.Wallhack.Float32() != 0 {
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.ColorBufferBit)
	}
	hooks.CallFuncInts0(rClearPattern.Address())
}

// HookedRDrawSequentialPoly R_DrawSequentialPoly
//export HookedRDrawSequentialPoly
func HookedRDrawSequentialPoly(surf uintptr, free int) {
	if cvar.Wallhack.Float32() == 0 {
		hooks.CallFuncInts2(rDrawSequentialPolyPattern.Address(), surf, uintptr(free))
		return
	}

	gl.Enable(gl.Blend)
	gl.DepthMask(false)
	gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
	gl.Color4f(1, 1, 1, cvar.WallhackAlpha.Float32())

	hooks.CallFuncInts2(rDrawSequentialPolyPattern.Address(), surf, uintptr(free))

	gl.DepthMask(true)
	gl.Disable(gl.Blend)
}

// HookedMemoryInit Memory_Init
//export HookedMemoryInit
func HookedMemoryInit(buf uintptr, size int) {
	hooks.CallFuncInts2(memoryInitPattern.Address(), buf, uintptr(size))
	registerCVars()
}

// TriGLRenderMode tri_GL_RenderMode
func TriGLRenderMode(mode int) {
	hooks.CallFuncInts1(triGLRenderModePattern.Address(), uintptr(mode))
}

// TriGLBegin tri_GL_Begin
func TriGLBegin(primitive int) {
	hooks.CallFuncInts1(triGLBeginPattern.Address(), uintptr(primitive))
}

// TriGLEnd tri_GL_End
func TriGLEnd() {
	hooks.CallFuncInts0(triGLEndPattern.Address())
}

// TriGLColor4f tri_GL_Color4f
func TriGLColor4f(r, g, b, a float32) {
	hooks.CallFuncFloats4(triGLColor4fPattern.Address(), r, g, b, a)
}

// TriGLVertex3fv tri_GL_Vertex3fv
func TriGLVertex3fv(v [3]float32) {
	hooks.CallFuncInts1(triGLVertex3fvPattern.Address(), uintptr(unsafe.Pointer(&v[0])))
}

// TriGLCullFace tri_GL_CullFace
func TriGLCullFace(style int) {
	hooks.CallFuncInts1(triGLCullFacePattern.Address(), uintptr(style))
}

// ScreenTransform ScreenTransform, similar to WorldToScreen in TriAPI
func ScreenTransform(point [3]float32) (screen [3]float32, clipped bool) {
	clipped = hooks.CallFuncInts2(screenTransformPattern.Address(), uintptr(unsafe.Pointer(&point[0])), uintptr(unsafe.Pointer(&screen[0]))) != 0
	return
}

// PFCheckClientI calls PF_checkclient_I, or FindClientInPVS as exported to DLL
func PFCheckClientI(edict unsafe.Pointer) uintptr {
	return uintptr(hooks.CallFuncInts1(pfCheckClientIPattern.Address(), uintptr(edict)))
}

// InitHWDLL initialise hw.dll hooks and symbol search
func InitHWDLL(base string) (err error) {
	if hwDLL != nil {
		return
	}

	hwDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&buildNumberPattern:                nil,
		&cvarRegisterVariablePattern:       nil,
		&vFadeAlphaPattern:                 C.HookedVFadeAlpha,
		&drawStringPattern:                 nil,
		&vgui2DrawSetTextColorAlphaPattern: nil,
		&hostAutoSaveFPattern:              nil,
		&hostNoclipFPattern:                nil,
		&pfTracelineDLLPattern:             nil,
		&triGLRenderModePattern:            nil,
		&triGLBeginPattern:                 nil,
		&triGLEndPattern:                   nil,
		&triGLColor4fPattern:               nil,
		&triGLCullFacePattern:              nil,
		&triGLVertex3fvPattern:             nil,
		&screenTransformPattern:            nil,
		&worldTransformPattern:             nil,
		&hudGetScreenInfoPattern:           nil,
		&rClearPattern:                     C.HookedRClear,
		&rDrawSequentialPolyPattern:        C.HookedRDrawSequentialPoly,
		&memoryInitPattern:                 C.HookedMemoryInit,
		&pfCheckClientIPattern:             nil,
	}

	errors := hooks.BatchFind(hwDLL, items)
	gamelibs.PrintBatchFindErrors(errors)

	switch hostAutoSaveFPattern.PatternKey() {
	case gamelibs.HL8684:
		address := *(*uintptr)(unsafe.Pointer(uintptr(hostAutoSaveFPattern.Address()) + 19))
		engine.Engine.SetSV(address)
		logs.DLLLog.Debugf("Set SV address: %x", address)
	}

	switch hostNoclipFPattern.PatternKey() {
	case gamelibs.HL8684:
		address := *(*uintptr)(unsafe.Pointer(uintptr(hostNoclipFPattern.Address()) + 31))
		address -= 0x14
		engine.Engine.SetGlobalVariables(address)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", address)
	}

	return
}
