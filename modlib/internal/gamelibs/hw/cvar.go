package hw

import (
	"hlinspect/internal/cvar"
)

func registerCVars() {
	CvarRegisterVariable(uintptr(cvar.Wallhack.Pointer()))
	CvarRegisterVariable(uintptr(cvar.WallhackAlpha.Pointer()))
	CvarRegisterVariable(uintptr(cvar.FadeRemove.Pointer()))
	CvarRegisterVariable(uintptr(cvar.Nodes.Pointer()))
	CvarRegisterVariable(uintptr(cvar.NodeLinks.Pointer()))
	CvarRegisterVariable(uintptr(cvar.Cine.Pointer()))
	CvarRegisterVariable(uintptr(cvar.CinePossess.Pointer()))
	CvarRegisterVariable(uintptr(cvar.CollisionHullShow.Pointer()))
}
