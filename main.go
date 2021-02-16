package main

import (
	"hlinspect/internal/gamelibs/cl"
	"hlinspect/internal/gamelibs/hl"
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
	defer func() {
		if r := recover(); r != nil {
			logs.DLLLog.Panicf("Got a panic during initialisation: %v", r)
		}
	}()

	logs.DLLLog.Debug("Initialising hooks")
	if !hooks.InitHooks() {
		logs.DLLLog.Panic("Unable to initialise hooks")
	}

	logs.DLLLog.Debug("Initialising HWDLL")
	if err := hw.InitHWDLL(); err != nil {
		logs.DLLLog.Panicf("Unable to initialise HWDLL: %v", err)
	}

	logs.DLLLog.Debug("Initialising ClientDLL")
	if err := cl.InitClientDLL(); err != nil {
		logs.DLLLog.Panicf("Unable to initialise ClientDLL: %v", err)
	}

	logs.DLLLog.Debug("Initialising HLDLL")
	if err := hl.InitHLDLL(); err != nil {
		logs.DLLLog.Panicf("Unable to initialise HLDLL: %v", err)
	}
}

// OnProcessDetach called from DllMain on process detach
//export OnProcessDetach
func OnProcessDetach() {
	hooks.CleanupHooks()
}
