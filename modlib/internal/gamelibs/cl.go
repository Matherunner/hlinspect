package gamelibs

import (
	"hlinspect/internal/gamelibs/cdefs"
	"hlinspect/internal/hooks"
	"unsafe"
)

var clientDLL *hooks.Module

func initAPIRegistry(reg *APIRegistry) {
	reg.HUDRedraw = hooks.MakeFunctionPattern("HUD_Redraw", map[string]string{"Windows": "HUD_Redraw"}, nil)
	reg.HUDDrawTransparentTriangles = hooks.MakeFunctionPattern("HUD_DrawTransparentTriangles", map[string]string{"Windows": "HUD_DrawTransparentTriangles"}, nil)
	reg.HUDVidInit = hooks.MakeFunctionPattern("HUD_VidInit", map[string]string{"Windows": "HUD_VidInit"}, nil)
	reg.HUDReset = hooks.MakeFunctionPattern("HUD_Reset", map[string]string{"Windows": "HUD_Reset"}, nil)
}

func initCLDLL(base string) (err error) {
	if clientDLL != nil {
		return
	}

	clientDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	reg := Model.Registry()

	initAPIRegistry(reg)

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&reg.HUDRedraw:                   cdefs.CDefs.HookedHUDRedraw,
		&reg.HUDDrawTransparentTriangles: cdefs.CDefs.HookedHUDDrawTransparentTriangles,
		&reg.HUDVidInit:                  cdefs.CDefs.HookedHUDVidInit,
		&reg.HUDReset:                    cdefs.CDefs.HookedHUDReset,
	}

	errors := hooks.BatchFind(clientDLL, items)
	printBatchFindErrors(errors)

	return nil
}
