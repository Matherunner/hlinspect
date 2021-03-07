package hw

import (
	"hlinspect/internal/cvar"
)

func registerCVars() {
	CvarRegisterVariable(uintptr(cvar.Wallhack.Pointer()))
	CvarRegisterVariable(uintptr(cvar.WallhackAlpha.Pointer()))
	CvarRegisterVariable(uintptr(cvar.FadeRemove.Pointer()))
}
