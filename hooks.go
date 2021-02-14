package main

/*
#cgo 386 LDFLAGS: -Llib -lMinHook

#include "MinHook.h"
*/
import "C"

func initHooks() {
	dllLog.Debug("Initializing hook")
	if err := C.MH_Initialize(); err != C.MH_OK {
		dllLog.Panicf("Unable to initialise hook: %v", err)
	}

}

func cleanupHooks() {
	dllLog.Debug("Uninitializing hook")
	C.MH_Uninitialize()
}
