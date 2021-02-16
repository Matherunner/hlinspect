package hooks

import (
	"fmt"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#cgo 386 LDFLAGS: -ldbghelp

#include <stdlib.h>
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <dbghelp.h>
*/
import "C"

var symLock sync.Mutex

func initSym() (ok bool) {
	procHandle := C.HANDLE(windows.CurrentProcess())
	if ret := C.SymInitialize(procHandle, nil, 1); ret == 0 {
		return
	}
	ok = true
	return
}

func cleanupSym() {
	procHandle := C.HANDLE(windows.CurrentProcess())
	C.SymCleanup(procHandle)
}

func findSym(name string) (relAddr uintptr, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var info C.SYMBOL_INFO_PACKAGE
	info.si.SizeOfStruct = C.ulong(unsafe.Sizeof(info.si))
	info.si.MaxNameLen = C.ulong(unsafe.Sizeof(info.name))

	procHandle := C.HANDLE(windows.CurrentProcess())
	symLock.Lock()
	ret := C.SymFromName(procHandle, cname, (*C.SYMBOL_INFO)(unsafe.Pointer(&info)))
	symLock.Unlock()
	if ret == 0 {
		lastErr := windows.GetLastError()
		err = fmt.Errorf("Unable to find symbol: %v", lastErr)
		return
	}

	relAddr = uintptr(info.si.Address - info.si.ModBase)
	return
}
