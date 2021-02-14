package cl

import (
	"hlinspect/internal/hooks"
	"hlinspect/internal/hud"
	"hlinspect/internal/logs"
)

/*
#include "defs.h"
*/
import "C"

var clientDLL *hooks.Module

var hudRedrawPattern = hooks.MakeFunctionPattern("HUD_Redraw", map[string]string{"Windows": "HUD_Redraw"}, nil)

func hudRedraw(time float32, intermission int32) {
	hooks.CallFuncFloatInt(hudRedrawPattern.Address(), time, uintptr(intermission))
}

// HookedHUDRedraw hooked HUD_Redraw
//export HookedHUDRedraw
func HookedHUDRedraw(time float32, intermission int32) {
	hudRedraw(time, intermission)
	hud.Draw(time, intermission)
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

	name, addr, err := hudRedrawPattern.Hook(clientDLL, C.HookedHUDRedraw)
	logs.DLLLog.Debugf("Found %v at %v using %v", hudRedrawPattern.Name(), addr, name)

	return
}
