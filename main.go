package main

import (
	"hlinspect/internal/gamelibs/hw"
	"hlinspect/internal/hooks"
)

/*
#include "dllmain.h"
*/
import "C"

func main() {}

// OnProcessAttach called from DllMain on process attach
//export OnProcessAttach
func OnProcessAttach() {
	hooks.InitHooks()
	hw.InitHWDLL()
}

// OnProcessDetach called from DllMain on process detach
//export OnProcessDetach
func OnProcessDetach() {
	hooks.CleanupHooks()
}
