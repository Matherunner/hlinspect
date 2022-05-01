package hooks

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

/*
#cgo 386 LDFLAGS: -ldbghelp

#include <stdlib.h>
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <dbghelp.h>
*/
import "C"

var (
	ErrCannotFindSymbol = errors.New("unable to find symbol")
)

var symLock sync.Mutex

// RefreshModuleList calls SymRefreshModuleList, must be called on DLL load
func RefreshModuleList() {
	symLock.Lock()
	defer symLock.Unlock()

	C.SymRefreshModuleList(C.GetCurrentProcess())
}

func initSym() (ok bool) {
	symLock.Lock()
	defer symLock.Unlock()

	if ret := C.SymInitialize(C.GetCurrentProcess(), nil, 1); ret == 0 {
		return
	}
	ok = true
	return
}

func cleanupSym() {
	symLock.Lock()
	defer symLock.Unlock()

	C.SymCleanup(C.GetCurrentProcess())
}

func findSym(name string) (relAddr uintptr, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var info C.SYMBOL_INFO_PACKAGE
	info.si.SizeOfStruct = C.ulong(unsafe.Sizeof(info.si))
	info.si.MaxNameLen = C.ulong(unsafe.Sizeof(info.name))

	symLock.Lock()
	defer symLock.Unlock()

	ret := C.SymFromName(C.GetCurrentProcess(), cname, (*C.SYMBOL_INFO)(unsafe.Pointer(&info)))
	lastErr := C.GetLastError()
	if ret == 0 {
		err = fmt.Errorf("%w: %v", ErrCannotFindSymbol, lastErr)
		return
	}

	relAddr = uintptr(info.si.Address - info.si.ModBase)
	return
}
