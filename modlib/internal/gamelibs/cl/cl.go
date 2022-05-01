package cl

import (
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gamelibs/common"
	"hlinspect/internal/hooks"
	"unsafe"
)

/*
#include "defs.h"
*/
import "C"

var clientDLL *hooks.Module

// HookedHUDRedraw hooked HUD_Redraw
//export HookedHUDRedraw
func HookedHUDRedraw(time float32, intermission int32) {
	gamelibs.Model.EventHandler().HUDRedraw(time, int(intermission))
}

// HookedHUDDrawTransparentTriangles HUD_DrawTransparentTriangles
//export HookedHUDDrawTransparentTriangles
func HookedHUDDrawTransparentTriangles() {
	gamelibs.Model.EventHandler().HUDDrawTransparentTriangle()
}

// HookedHUDVidInit HUD_VidInit
//export HookedHUDVidInit
func HookedHUDVidInit() int {
	return gamelibs.Model.EventHandler().HUDVidInit()
}

// HookedHUDReset HUD_Reset
//export HookedHUDReset
func HookedHUDReset() {
	gamelibs.Model.EventHandler().HUDReset()
}

func initAPIRegistry(reg *gamelibs.APIRegistry) {
	reg.HUDRedraw = hooks.MakeFunctionPattern("HUD_Redraw", map[string]string{"Windows": "HUD_Redraw"}, nil)
	reg.HUDDrawTransparentTriangles = hooks.MakeFunctionPattern("HUD_DrawTransparentTriangles", map[string]string{"Windows": "HUD_DrawTransparentTriangles"}, nil)
	reg.HUDVidInit = hooks.MakeFunctionPattern("HUD_VidInit", map[string]string{"Windows": "HUD_VidInit"}, nil)
	reg.HUDReset = hooks.MakeFunctionPattern("HUD_Reset", map[string]string{"Windows": "HUD_Reset"}, nil)
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

	reg := gamelibs.Model.Registry()

	initAPIRegistry(reg)

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&reg.HUDRedraw:                   C.HookedHUDRedraw,
		&reg.HUDDrawTransparentTriangles: C.HookedHUDDrawTransparentTriangles,
		&reg.HUDVidInit:                  C.HookedHUDVidInit,
		&reg.HUDReset:                    C.HookedHUDReset,
	}

	errors := hooks.BatchFind(clientDLL, items)
	common.PrintBatchFindErrors(errors)

	return
}
