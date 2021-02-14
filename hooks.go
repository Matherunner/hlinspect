package main

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#cgo 386 LDFLAGS: -Llib -lMinHook

#include <windows.h>
#include "MinHook.h"
*/
import "C"

func getModuleInfo(name string) (base uintptr, size uint, err error) {
	type moduleInfo struct {
		BaseOfDll   uintptr
		SizeOfImage uint32
		EntryPoint  uintptr
	}

	cname, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return
	}

	var handle windows.Handle
	err = windows.GetModuleHandleEx(0, cname, &handle)
	if err != nil {
		return
	}
	defer windows.FreeLibrary(handle)

	psapi := windows.MustLoadDLL("psapi.dll")
	getModuleInformation := psapi.MustFindProc("GetModuleInformation")
	procHandle := windows.CurrentProcess()

	var info moduleInfo
	ret, _, _ := getModuleInformation.Call(uintptr(procHandle), uintptr(handle), uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if ret == 0 {
		err = errors.New("GetModuleInformation failed")
		return
	}

	base = info.BaseOfDll
	size = uint(info.SizeOfImage)
	return
}

func initHooks() {
	base, size, err := getModuleInfo("hw.dll")
	if err != nil {
		dllLog.Panic("Unable to get module info")
	}

	dllLog.Debugf("Got base and size of hw.dll: %v %v\n", base, size)

	dllLog.Debug("Initializing hook")
	if err := C.MH_Initialize(); err != C.MH_OK {
		dllLog.Panicf("Unable to initialise hook: %v", err)
	}

}

func cleanupHooks() {
	dllLog.Debug("Uninitializing hook")
	C.MH_Uninitialize()
}
