package cl

import (
	"hlinspect/internal/graphics"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"unsafe"
)

/*
#cgo 386 LDFLAGS: -lopengl32

#include <GL/gl.h>
#include "defs.h"
*/
import "C"

var clientDLL *hooks.Module

var hudRedrawPattern = hooks.MakeFunctionPattern("HUD_Redraw", map[string]string{"Windows": "HUD_Redraw"}, nil)
var hudDrawTransparentTriangles = hooks.MakeFunctionPattern("HUD_DrawTransparentTriangles", map[string]string{"Windows": "HUD_DrawTransparentTriangles"}, nil)

// HookedHUDRedraw hooked HUD_Redraw
//export HookedHUDRedraw
func HookedHUDRedraw(time float32, intermission int32) {
	hooks.CallFuncFloatInt(hudRedrawPattern.Address(), time, uintptr(intermission))
	graphics.DrawHUD(time, intermission)
}

// HookedHUDDrawTransparentTriangles HUD_DrawTransparentTriangles
//export HookedHUDDrawTransparentTriangles
func HookedHUDDrawTransparentTriangles() {
	hooks.CallFuncInts0(hudDrawTransparentTriangles.Address())
	C.glDisable(C.GL_TEXTURE_2D)
	graphics.DrawTriangles()
	C.glEnable(C.GL_TEXTURE_2D)
}

// InitClientDLL initialise client.dll
func InitClientDLL() (err error) {
	if clientDLL != nil {
		return
	}

	clientDLL, err = hooks.NewModule("client.dll")
	if err != nil {
		return
	}

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&hudRedrawPattern:            C.HookedHUDRedraw,
		&hudDrawTransparentTriangles: C.HookedHUDDrawTransparentTriangles,
	}

	errors := hooks.BatchFind(clientDLL, items)
	for pat, err := range errors {
		if err == nil {
			logs.DLLLog.Debugf("Found %v at %v", pat.Name(), pat.Address())
		} else {
			logs.DLLLog.Debugf("Failed to find %v: %v", pat.Name(), err)
		}
	}

	return
}
