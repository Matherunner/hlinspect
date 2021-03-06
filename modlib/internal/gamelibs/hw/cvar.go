package hw

import (
	"hlinspect/internal/cvar"
)

func registerCVars() {
	CvarRegisterVariable(uintptr(cvar.Wallhack.Pointer()))
}
