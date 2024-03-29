package game

import (
	"hlinspect/internal/game/cdefs"
	"hlinspect/internal/hooks"
	"unsafe"
)

var clientDLL *hooks.Module

func initCLDLL(base string) (err error) {
	if clientDLL != nil {
		return
	}

	clientDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	reg := Model.Registry()

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&reg.CLCreateMove:                cdefs.CDefs.HookedCLCreateMove,
		&reg.HUDRedraw:                   cdefs.CDefs.HookedHUDRedraw,
		&reg.HUDDrawTransparentTriangles: cdefs.CDefs.HookedHUDDrawTransparentTriangles,
		&reg.HUDVidInit:                  cdefs.CDefs.HookedHUDVidInit,
		&reg.HUDReset:                    cdefs.CDefs.HookedHUDReset,
	}

	errors := hooks.BatchFind(clientDLL, items)
	printBatchFindErrors(errors)

	return nil
}
