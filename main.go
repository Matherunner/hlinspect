package main

import (
	"hlinspect/internal/gamelibs/hw"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
)

/*
#include "dllmain.h"
*/
import "C"

func main() {}

// OnProcessAttach called from DllMain on process attach
//export OnProcessAttach
func OnProcessAttach() {
	logs.DLLLog.Debug("Initialising hooks")
	if !hooks.InitHooks() {
		logs.DLLLog.Panic("Unable to initialise hooks")
	}

	logs.DLLLog.Debug("Initialising HWDLL")
	if err := hw.InitHWDLL(); err != nil {
		logs.DLLLog.Panicf("Unable to initialise HWDLL: %v", err)
	}
}

// OnProcessDetach called from DllMain on process detach
//export OnProcessDetach
func OnProcessDetach() {
	hooks.CleanupHooks()
}
