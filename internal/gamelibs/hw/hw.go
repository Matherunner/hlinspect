package hw

import (
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"strconv"
	"unsafe"
)

/*
#include <stdlib.h>
#include "defs.h"
*/
import "C"

var hwDLL *hooks.Module

var buildNumberPattern = hooks.MakeFunctionPattern("build_number", nil, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
})
var cvarRegisterVariablePattern = hooks.MakeFunctionPattern("Cvar_RegisterVariable", nil, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC 83 EC 14 53 56 8B 75 08 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68"),
	"HL-NGHL":      hooks.MustMakePattern("83 EC 14 53 56 8B 74 24 20 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5B 83 C4 14 C3 8B 16 52 E8"),
})
var vFadeAlphaPattern = hooks.MakeFunctionPattern("V_FadeAlpha", nil, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC 83 EC 08 D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? DF E0 F6 C4 05 7A 1C D9 05 ?? ?? ?? ?? DC 1D"),
	"HL-4554":      hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? 83 EC"),
})
var drawStringPattern = hooks.MakeFunctionPattern("Draw_String", nil, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC 56 57 E8 ?? ?? ?? ?? 8B 4D 0C 8B 75 08 50 8B 45 10 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
})
var vgui2DrawSetTextColorAlphaPattern = hooks.MakeFunctionPattern("VGUI2_Draw_SetTextColorAlpha", nil, map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC 8A 45 08 8A 4D 0C 8A 55 10 88 45 08 8A 45 14 88 4D 09 88 55 0A 88 45 0B 8B 4D 08 89"),
})

// BuildNumber build_number
func BuildNumber() int {
	return hooks.CallFuncInts0(buildNumberPattern.Address())
}

// CvarRegisterVariable Cvar_RegisterVariable
func CvarRegisterVariable(name string, value string) {
	floatVal, _ := strconv.ParseFloat(value, 32)
	cvar := rawCVar{
		// Probably don't need to free these strings?
		Name:   uintptr(unsafe.Pointer(C.CString(name))),
		String: uintptr(unsafe.Pointer(C.CString(value))),
		Value:  float32(floatVal),
	}
	hooks.CallFuncInts1(cvarRegisterVariablePattern.Address(), uintptr(unsafe.Pointer(&cvar)))
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
	return 0
}

// InitHWDLL initialise hw.dll hooks and symbol search
func InitHWDLL() (err error) {
	if hwDLL != nil {
		return
	}

	hwDLL, err = hooks.NewModule("hw.dll")
	if err != nil {
		return
	}

	var name string
	var addr unsafe.Pointer

	name, addr, err = buildNumberPattern.Find(hwDLL)
	logs.DLLLog.Debugf("Found %v at %v using %v", buildNumberPattern.Name(), addr, name)

	name, addr, err = cvarRegisterVariablePattern.Find(hwDLL)
	logs.DLLLog.Debugf("Found %v at %v using %v", cvarRegisterVariablePattern.Name(), addr, name)

	name, addr, err = vFadeAlphaPattern.Hook(hwDLL, C.HookedVFadeAlpha)
	logs.DLLLog.Debugf("Found %v at %v using %v", vFadeAlphaPattern.Name(), addr, name)

	name, addr, err = drawStringPattern.Find(hwDLL)
	logs.DLLLog.Debugf("Found %v at %v using %v", drawStringPattern.Name(), addr, name)

	name, addr, err = vgui2DrawSetTextColorAlphaPattern.Find(hwDLL)
	logs.DLLLog.Debugf("Found %v at %v using %v", vgui2DrawSetTextColorAlphaPattern.Name(), addr, name)

	return
}
