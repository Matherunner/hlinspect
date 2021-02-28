package cl

import (
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gamelibs/hw"
	"hlinspect/internal/graphics"
	"hlinspect/internal/hooks"
	"unsafe"
)

/*
#include "defs.h"
*/
import "C"

var clientDLL *hooks.Module

var hudRedrawPattern = hooks.MakeFunctionPattern("HUD_Redraw", map[string]string{"Windows": "HUD_Redraw"}, nil)
var hudDrawTransparentTrianglesPattern = hooks.MakeFunctionPattern("HUD_DrawTransparentTriangles", map[string]string{"Windows": "HUD_DrawTransparentTriangles"}, nil)
var hudVidInitPattern = hooks.MakeFunctionPattern("HUD_VidInit", map[string]string{"Windows": "HUD_VidInit"}, nil)
var hudResetPattern = hooks.MakeFunctionPattern("HUD_Reset", map[string]string{"Windows": "HUD_Reset"}, nil)

// HookedHUDRedraw hooked HUD_Redraw
//export HookedHUDRedraw
func HookedHUDRedraw(time float32, intermission int32) {
	hooks.CallFuncFloatInt(hudRedrawPattern.Address(), time, uintptr(intermission))
	graphics.DrawHUD(time, intermission)
}

// HookedHUDDrawTransparentTriangles HUD_DrawTransparentTriangles
//export HookedHUDDrawTransparentTriangles
func HookedHUDDrawTransparentTriangles() {
	hooks.CallFuncInts0(hudDrawTransparentTrianglesPattern.Address())
	graphics.GLDisable(graphics.GLTexture2D)
	graphics.DrawTriangles()
	graphics.GLEnable(graphics.GLTexture2D)
}

// HookedHUDVidInit HUD_VidInit
//export HookedHUDVidInit
func HookedHUDVidInit() int {
	ret := hooks.CallFuncInts0(hudVidInitPattern.Address())
	screenInfo := hw.GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)
	return ret
}

// HookedHUDReset HUD_Reset
//export HookedHUDReset
func HookedHUDReset() {
	hooks.CallFuncInts0(hudResetPattern.Address())
	screenInfo := hw.GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)
}

// InitClientDLL initialise client.dll
func InitClientDLL(base string) (err error) {
	if clientDLL != nil {
		return
	}

	clientDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&hudRedrawPattern:                   C.HookedHUDRedraw,
		&hudDrawTransparentTrianglesPattern: C.HookedHUDDrawTransparentTriangles,
		&hudVidInitPattern:                  C.HookedHUDVidInit,
		&hudResetPattern:                    C.HookedHUDReset,
	}

	errors := hooks.BatchFind(clientDLL, items)
	gamelibs.PrintBatchFindErrors(errors)

	return
}
