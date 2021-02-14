package main

import (
	"errors"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#cgo 386 LDFLAGS: -Llib -lMinHook

#include <stdlib.h>
#include <windows.h>
#include "MinHook.h"

typedef int (*args0Func)();

static inline int call_args_0(void *f)
{
	return ((args0Func)f)();
}

extern int buildNumberNew();
*/
import "C"

type originalAddr struct {
	buildNumber unsafe.Pointer
}

type hwDLL struct{}

type searchPattern struct {
	Bytes  []uint8
	Ignore []bool
}

type functionPattern struct {
	FunctionName string
	SymbolName   string
	Patterns     []searchPattern
	AddrPointer  *unsafe.Pointer
	ReplaceAddr  unsafe.Pointer
}

func (hw *hwDLL) BuildNumber() int {
	return -1
}

//export buildNumberNew
func buildNumberNew() int {
	return 1337
}

var originalAddress originalAddr

var patterns = []functionPattern{
	functionPattern{
		FunctionName: "build_number",
		Patterns: []searchPattern{
			mustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
		},
		AddrPointer: &originalAddress.buildNumber,
	},
}

func mustMakePattern(pattern string) searchPattern {
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
	return searchPattern{Bytes: patternBytes, Ignore: ignoreBytes}
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

func findSubstringPattern(pattern searchPattern, base unsafe.Pointer, size uint) (addr unsafe.Pointer, found bool, err error) {
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

func initHooks() {
	base, size, err := getModuleInfo("hw.dll")
	if err != nil {
		dllLog.Panicf("Unable to get module info: %v", err)
	}
	dllLog.Debugf("Got base and size of hw.dll: %v %v\n", base, size)

	dllLog.Debugf("before pattern search: %v\n", patterns)
	addr, found, err := findSubstringPattern(patterns[0].Patterns[0], base, size)
	dllLog.Debugf("pattern search: %v %v %v\n", addr, found, err)

	dllLog.Debug("Initializing hook")
	if err := C.MH_Initialize(); err != C.MH_OK {
		dllLog.Panicf("Unable to initialise hook: %v", err)
	}

	dllLog.Debug("Hooking")
	var origFunc uintptr
	if err := C.MH_CreateHook(C.LPVOID(addr), C.LPVOID(C.buildNumberNew), (*C.LPVOID)(unsafe.Pointer(&origFunc))); err != C.MH_OK {
		dllLog.Debugf("Unable to hook build number!")
	}

	if err := C.MH_EnableHook(C.LPVOID(nil)); err != C.MH_OK {
		dllLog.Debugf("Unable to enable hook")
	}
}

func cleanupHooks() {
	dllLog.Debug("Uninitializing hook")
	C.MH_Uninitialize()
}
