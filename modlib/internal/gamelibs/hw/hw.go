package hw

import (
	"context"
	"hlinspect/internal/cvar"
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gamelibs/common"
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
	common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
	common.HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 83 EC 08 57 33 FF 85 C0 0F 85 A5 00 00 00 53 56 33 DB BE ?? ?? ?? ?? 8B 06 8B 0D"),
	common.HLNGHL: hooks.MustMakePattern("A1 ?? ?? ?? ?? 83 EC 08 56 33 F6 85 C0 0F 85 9F 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D"),
})
var vFadeAlphaPattern = hooks.MakeFunctionPattern("V_FadeAlpha", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? DF E0 F6 C4 05 7A 1C D9 05 ?? ?? ?? ?? DC 1D"),
	common.HL4554: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? 83 EC 08 DF E0 F6 C4"),
})
var drawStringPattern = hooks.MakeFunctionPattern("Draw_String", nil, map[string]hooks.SearchPattern{
	// Search for "%i %i %i", there is a thunk that call this function. There are two functions that call the thunk.
	// The shorter one is Draw_String.
	common.HL8684: hooks.MustMakePattern("55 8B EC 56 57 E8 ?? ?? ?? ?? 8B 4D 0C 8B 75 08 50 8B 45 10 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
	common.HL4554: hooks.MustMakePattern("56 57 E8 ?? ?? ?? ?? 8B 4C 24 10 8B 74 24 0C 50 8B 44 24 18 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
})
var vgui2DrawSetTextColorAlphaPattern = hooks.MakeFunctionPattern("VGUI2_Draw_SetTextColorAlpha", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 8A 45 08 8A 4D 0C 8A 55 10 88 45 08 8A 45 14 88 4D 09 88 55 0A 88 45 0B 8B 4D 08 89"),
	common.HL4554: hooks.MustMakePattern("8A 44 24 04 8A 4C 24 08 8A 54 24 0C 88 44 24 04 8A 44 24 10 88 4C 24 05 88 54 24 06 88 44 24 07 8B 4C 24 04 89 0D"),
})
var hostAutoSaveFPattern = hooks.MakeFunctionPattern("Host_AutoSave_f", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? B9 01 00 00 00 3B C1 0F 85 9F 00 00 00 A1 ?? ?? ?? ?? 85 C0 75 10 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 04 33 C0 C3 39 0D"),
})
var hostNoclipFPattern = hooks.MakeFunctionPattern("Host_Noclip_f", nil, map[string]hooks.SearchPattern{
	// Search for "noclip ON\n"
	common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 24 A1 ?? ?? ?? ?? BA 01 00 00 00 3B C2 75 09 E8 ?? ?? ?? ?? 8B E5 5D C3 D9 05 ?? ?? ?? ?? D8 1D"),
	common.HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? BA 01 00 00 00 83 EC 24 3B C2 75 09 E8 ?? ?? ?? ?? 83 C4 24 C3 D9 05 ?? ?? ?? ?? D8 1D"),
	common.HLNGHL: hooks.MustMakePattern("A1 ?? ?? ?? ?? BA 01 00 00 00 83 EC 24 3B C2 75 08 83 C4 24 E9 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? D8 1D"),
})
var triGLRenderModePattern = hooks.MakeFunctionPattern("tri_GL_RenderMode", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 56 8B 75 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
	common.HL4554: hooks.MustMakePattern("56 8B 74 24 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
})
var triGLBeginPattern = hooks.MakeFunctionPattern("tri_GL_Begin", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 45 08 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? 5D C3"),
	common.HL4554: hooks.MustMakePattern("E8 ?? ?? ?? ?? 8B 44 24 04 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? C3"),
})
var triGLEndPattern = hooks.MakeFunctionPattern("tri_GL_End", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 55 8B EC 8B 45 0C"),
	common.HL4554: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 8B 44 24 08 8B 4C 24 04"),
})
var triGLColor4fPattern = hooks.MakeFunctionPattern("tri_GL_Color4f", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 51 83 3D ?? ?? ?? ?? 04 75 4A D9 45 14 D8 0D ?? ?? ?? ?? D9 5D FC D9 45 FC E8 ?? ?? ?? ?? D9 45 10"),
	common.HL4554: hooks.MustMakePattern("51 83 3D ?? ?? ?? ?? 04 75 50 D9 44 24 14 D8 0D ?? ?? ?? ?? D9 5C 24 00 D9 44 24 00 E8 ?? ?? ?? ?? D9 44 24 10"),
})
var triGLCullFacePattern = hooks.MakeFunctionPattern("tri_GL_CullFace", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 83 E8 00 74 10 48 75 23 68 44 0B 00 00 FF 15 ?? ?? ?? ?? 5D C3 68 44 0B 00 00"),
	common.HL4554: hooks.MustMakePattern("8B 44 24 04 83 E8 00 74 0F 48 75 22 68 44 0B 00 00 FF 15 ?? ?? ?? ?? C3 68 44 0B 00 00"),
})
var triGLVertex3fvPattern = hooks.MakeFunctionPattern("tri_GL_Vertex3fv", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 50 FF 15 ?? ?? ?? ?? 5D C3 90 55 8B EC 8B 45 10 8B 4D 0C 8B 55 08 50 51 52"),
	common.HL4554: hooks.MustMakePattern("8B 44 24 04 50 FF 15 ?? ?? ?? ?? C3 90 90 90 90 8B 44 24 0C 8B 4C 24 08 8B 54 24 04 50 51 52"),
})
var screenTransformPattern = hooks.MakeFunctionPattern("ScreenTransform", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 51 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
	common.HL4554: hooks.MustMakePattern("51 8B 44 24 08 8B 4C 24 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
})
var worldTransformPattern = hooks.MakeFunctionPattern("WorldTransform", nil, map[string]hooks.SearchPattern{
	// Most likely the function below ScreenTransform
	common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48"),
	common.HL4554: hooks.MustMakePattern("83 EC 08 8B 44 24 0C 8B 4C 24 10 D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08"),
})
var hudGetScreenInfoPattern = hooks.MakeFunctionPattern("hudGetScreenInfo", nil, map[string]hooks.SearchPattern{
	// Search for "Half-Life %i/%s (hw build %d)". This function is Draw_ConsoleBackground
	// The function below it should be Draw_FillRGBA. Get the cross references to Draw_FillRGBA. One of them
	// is a global variable of enginefuncs. The next entry is hudGetScreenInfo.
	common.HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 83 C4 04 85 C0 75 02 5D C3 81 38 14 02 00 00 74 04"),
	common.HL4554: hooks.MustMakePattern("8D 44 24 04 50 FF 15 ?? ?? ?? ?? 8B 44 24 08 83 C4 04 85 C0 75 01 C3 81 38 14 02 00 00 74 03"),
})
var rDrawSequentialPolyPattern = hooks.MakeFunctionPattern("R_DrawSequentialPoly", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 53 56 57 83 B8 F8 02 00 00 01 75 63 E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
	common.HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 53 55 56 8B 88 F8 02 00 00 BE 01 00 00 00 3B CE 57 75 ?? E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
})
var rClearPattern = hooks.MakeFunctionPattern("R_Clear", nil, map[string]hooks.SearchPattern{
	common.HL8684: hooks.MustMakePattern("8B 15 ?? ?? ?? ?? 33 C0 83 FA 01 0F 9F C0 50 E8 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 83 C4 04 DF E0"),
	common.HLNGHL: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 F6 C4 44 7B 34 D9 05 ?? ?? ?? ?? D8 1D"),
})
var pfCheckClientIPattern = hooks.MakeFunctionPattern("PF_checkclient_I", nil, map[string]hooks.SearchPattern{
	// Search for "Spawned a NULL entity!", the referencing function is CreateNamedEntity
	// Find cross references, go to the global data, that data is g_engfuncsExportedToDlls
	// Go up 6 entries, we will end up at PF_checkclient_I
	common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 0C DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
	common.HL4554: hooks.MustMakePattern("DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? 83 EC 0C DC 1D ?? ?? ?? ?? DF E0 F6 C4 01 A1 ?? ?? ?? ?? 75 26"),
	common.HLNGHL: hooks.MustMakePattern("DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? 83 EC 0C DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
})

// GetScreenInfo proxies hudGetScreenInfo
func GetScreenInfo() gamelibs.ScreenInfo {
	screenInfo := gamelibs.ScreenInfo{}
	screenInfo.Size = int32(unsafe.Sizeof(screenInfo))
	hooks.CallFuncInts1(hudGetScreenInfoPattern.Address(), uintptr(unsafe.Pointer(&screenInfo)))
	return screenInfo
}

// BuildNumber build_number
func BuildNumber() int {
	return hooks.CallFuncInts0(buildNumberPattern.Address())
}

// CmdHandler called by C code. This is needed because passing Go function directly to CmdAddCommand doesn't seem to work.
//export CmdHandler
func CmdHandler() {
	gamelibs.Model.EventHandler().OnCommand()
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
	gamelibs.Model.EventHandler().MemoryInit(context.TODO(), buf, size)
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

func initAPIRegistry(reg *gamelibs.APIRegistry) {
	reg.CvarRegisterVariable = hooks.MakeFunctionPattern("Cvar_RegisterVariable", nil, map[string]hooks.SearchPattern{
		common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 14 53 56 8B 75 08 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68"),
		common.HLNGHL: hooks.MustMakePattern("83 EC 14 53 56 8B 74 24 20 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5B 83 C4 14 C3 8B 16 52 E8"),
	})
	reg.MemoryInit = hooks.MakeFunctionPattern("Memory_Init", nil, map[string]hooks.SearchPattern{
		common.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 8B 4D 0C 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
		common.HL4554: hooks.MustMakePattern("8B 44 24 04 8B 4C 24 08 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
	})
	reg.CmdAddCommandWithFlags = hooks.MakeFunctionPattern("Cmd_AddCommandWithFlags", nil, map[string]hooks.SearchPattern{
		// Search for "Cmd_AddCommand: %s already defined as a var"
		common.HL8684: hooks.MustMakePattern("55 8B EC 56 57 8B 7D 08 57 E8 ?? ?? ?? ?? 8A 08 83 C4 04 84 C9 74 12 57 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5D C3 8B 35"),
		common.HLNGHL: hooks.MustMakePattern("56 57 8B 7C 24 0C 57 E8 ?? ?? ?? ?? 8A 08 83 C4 04 84 C9 74 11 57 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E C3 8B 35"),
	})
	reg.CmdArgv = hooks.MakeFunctionPattern("Cmd_Argv", nil, map[string]hooks.SearchPattern{
		// Search for "MISSING VALUE" to find Host_FullInfo_f
		// The first function called is Cmd_Argc, while the second function with one argument should be Cmd_Argv
		common.HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 8B 0D ?? ?? ?? ?? 83 C4 04 3B C1 72 07 A1 ?? ?? ?? ?? 5D"),
		common.HL4554: hooks.MustMakePattern("8D 44 24 04 50 FF 15 ?? ?? ?? ?? 8B 44 24 08 8B 0D ?? ?? ?? ?? 83 C4 04 3B C1 72 06 A1 ?? ?? ?? ?? C3"),
	})
	reg.AngleVectors = hooks.MakeFunctionPattern("AngleVectors", nil, map[string]hooks.SearchPattern{
		common.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 1C 8D 45 14 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 83 C4 08"),
		common.HL4554: hooks.MustMakePattern("55 8B EC 83 E4 F8 83 EC 20 56 8D 45 14 57 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 D9 41 04"),
		common.HLNGHL: hooks.MustMakePattern("55 8B EC 83 E4 F8 83 EC 20 8D 45 14 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 83 C4 08"),
	})
	reg.PFTracelineDLL = hooks.MakeFunctionPattern("PF_traceline_DLL", nil, map[string]hooks.SearchPattern{
		common.HL8684: hooks.MustMakePattern("55 8B EC 8B 45 14 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4D 0C 8B 55 08 56 50 8B 45 10 50 51 52 E8 ?? ?? ?? ?? D9 05"),
		common.HL4554: hooks.MustMakePattern("8B 44 24 10 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4C 24 08 8B 54 24 04 56 50 8B 44 24 14 50 51 52 E8 ?? ?? ?? ?? D9 05"),
	})

	reg.CCmdHandler = C.CCmdHandler
}

// InitHWDLL initialise hw.dll hooks and symbol search with idempotency.
func InitHWDLL(base string) (err error) {
	if hwDLL != nil {
		return
	}

	hwDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	reg := gamelibs.Model.Registry()

	initAPIRegistry(reg)

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&buildNumberPattern:                nil,
		&reg.CvarRegisterVariable:          nil,
		&reg.CmdAddCommandWithFlags:        nil,
		&reg.CmdArgv:                       nil,
		&vFadeAlphaPattern:                 C.HookedVFadeAlpha,
		&drawStringPattern:                 nil,
		&vgui2DrawSetTextColorAlphaPattern: nil,
		&hostAutoSaveFPattern:              nil,
		&hostNoclipFPattern:                nil,
		&reg.PFTracelineDLL:                nil,
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
		&reg.MemoryInit:                    C.HookedMemoryInit,
		&pfCheckClientIPattern:             nil,
		&reg.AngleVectors:                  nil,
	}

	errors := hooks.BatchFind(hwDLL, items)
	common.PrintBatchFindErrors(errors)

	switch hostAutoSaveFPattern.PatternKey() {
	case common.HL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(hostAutoSaveFPattern.Address()) + 19))
		engine.Engine.SetSV(ptr)
		logs.DLLLog.Debugf("Set SV address: %x", ptr)
	}

	switch hostNoclipFPattern.PatternKey() {
	case common.HL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(hostNoclipFPattern.Address()) + 31))
		ptr = unsafe.Pointer(uintptr(ptr) - 0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case common.HL4554:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(hostNoclipFPattern.Address()) + 28))
		ptr = unsafe.Pointer(uintptr(ptr) - 0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case common.HLNGHL:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(hostNoclipFPattern.Address()) + 27))
		ptr = unsafe.Pointer(uintptr(ptr) - 0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	}

	return
}
