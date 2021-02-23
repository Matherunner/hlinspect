package hw

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs"
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
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 56 57 E8 ?? ?? ?? ?? 8B 4D 0C 8B 75 08 50 8B 45 10 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
})
var vgui2DrawSetTextColorAlphaPattern = hooks.MakeFunctionPattern("VGUI2_Draw_SetTextColorAlpha", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 8A 45 08 8A 4D 0C 8A 55 10 88 45 08 8A 45 14 88 4D 09 88 55 0A 88 45 0B 8B 4D 08 89"),
})
var hostAutoSaveFPattern = hooks.MakeFunctionPattern("Host_AutoSave_f", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? B9 01 00 00 00 3B C1 0F 85 9F 00 00 00 A1 ?? ?? ?? ?? 85 C0 75 10 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 04 33 C0 C3 39 0D"),
})
var hostNoclipFPattern = hooks.MakeFunctionPattern("Host_Noclip_f", nil, map[string]hooks.SearchPattern{
	gamelibs.HL8684: hooks.MustMakePattern("55 8B EC 83 EC 24 A1 ?? ?? ?? ?? BA 01 00 00 00 3B C2 75 09 E8 ?? ?? ?? ?? 8B E5 5D C3 D9 05 ?? ?? ?? ?? D8 1D"),
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

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&buildNumberPattern:                nil,
		&cvarRegisterVariablePattern:       nil,
		&vFadeAlphaPattern:                 C.HookedVFadeAlpha,
		&drawStringPattern:                 nil,
		&vgui2DrawSetTextColorAlphaPattern: nil,
		&hostAutoSaveFPattern:              nil,
		&hostNoclipFPattern:                nil,
	}

	errors := hooks.BatchFind(hwDLL, items)
	for pat, err := range errors {
		if err == nil {
			logs.DLLLog.Debugf("Found %v at %v", pat.Name(), pat.Address())
		} else {
			logs.DLLLog.Debugf("Failed to find %v: %v", pat.Name(), err)
		}
	}

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
