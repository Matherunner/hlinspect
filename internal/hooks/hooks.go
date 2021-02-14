package hooks

import (
	"errors"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#cgo 386 LDFLAGS: -Llib -lMinHook

#include <windows.h>
#include "MinHook.h"
*/
import "C"

type Module struct {
	name string
	base unsafe.Pointer
	size uint
}

func NewModule(name string) (module *Module, err error) {
	base, size, err := getModuleInfo(name)
	if err != nil {
		return
	}

	module = &Module{name: name, base: base, size: size}
	return
}

type SearchPattern struct {
	Bytes  []uint8
	Ignore []bool
}

type FunctionPattern struct {
	functionName string
	symbolName   string
	patterns     map[string]SearchPattern
	addrPointer  unsafe.Pointer
	replaceAddr  unsafe.Pointer
}

func MakeFunctionPattern(functionName string, patterns map[string]SearchPattern) FunctionPattern {
	return FunctionPattern{functionName: functionName, patterns: patterns}
}

func (pat *FunctionPattern) Find(module *Module) (foundName string, address unsafe.Pointer, err error) {
	for patternName, pattern := range pat.patterns {
		addr, found, err := findSubstringPattern(pattern, module.base, module.size)
		if err != nil {
			break
		}
		if found {
			foundName = patternName
			address = addr
			pat.addrPointer = addr
			break
		}
	}
	return
}

func (pat *FunctionPattern) GetAddress() unsafe.Pointer {
	return pat.addrPointer
}

// MustMakePattern create a SearchPattern, panic if pattern is malformed
func MustMakePattern(pattern string) SearchPattern {
	patternTokens := strings.Split(pattern, " ")
	patternBytes := make([]uint8, len(patternTokens))
	ignoreBytes := make([]bool, len(patternTokens))
	for i, token := range patternTokens {
		if token == "??" {
			ignoreBytes[i] = true
			continue
		}
		val, err := strconv.ParseInt(token, 16, 0)
		if err != nil {
			panic("Invalid pattern given")
		}
		patternBytes[i] = uint8(val)
	}
	return SearchPattern{Bytes: patternBytes, Ignore: ignoreBytes}
}

func getModuleInfo(name string) (base unsafe.Pointer, size uint, err error) {
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
	base = unsafe.Pointer(info.BaseOfDll)
	size = uint(info.SizeOfImage)
	return
}

func findSubstringPattern(pattern SearchPattern, base unsafe.Pointer, size uint) (addr unsafe.Pointer, found bool, err error) {
	patternLen := uint(len(pattern.Bytes))
	for i := uint(0); i < size-patternLen; i++ {
		found = true
		for j := uint(0); j < patternLen; j++ {
			if pattern.Ignore[j] {
				continue
			}
			if *(*uint8)(unsafe.Pointer(uintptr(base) + uintptr(i+j))) != pattern.Bytes[j] {
				found = false
				break
			}
		}
		if found {
			addr = unsafe.Pointer(uintptr(base) + uintptr(i))
			break
		}
	}
	return
}

var minHookInitialized = false

// InitHooks initialise hooks
func InitHooks() (ok bool) {
	if minHookInitialized {
		ok = true
		return
	}
	if err := C.MH_Initialize(); err != C.MH_OK {
		return
	}
	minHookInitialized = true
	ok = true
	return
}

// CleanupHooks cleanup hooks
func CleanupHooks() {
	C.MH_Uninitialize()
}
