package hooks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#cgo 386 LDFLAGS: -Llib -lMinHook

#include "MinHook.h"
*/
import "C"

var (
	ErrPatternNotFound  = errors.New("pattern not found")
	ErrNonUniquePattern = errors.New("non-unique pattern")
	ErrCreateHookFailed = errors.New("failed to create hook")
	ErrEnableHookFailed = errors.New("failed to enable hook")
)

type PatternMap map[string]SearchPattern

type SymbolNameMap map[string]string

type Module struct {
	name   string
	base   unsafe.Pointer
	size   uint
	handle windows.Handle
}

func NewModule(name string) (module *Module, err error) {
	type moduleInfo struct {
		BaseOfDll   uintptr
		SizeOfImage uint32
		EntryPoint  uintptr
	}

	var cname *uint16
	cname, err = windows.UTF16PtrFromString(name)
	if err != nil {
		return
	}

	var handle windows.Handle
	err = windows.GetModuleHandleEx(0, cname, &handle)
	if err != nil {
		return
	}

	psapi := windows.MustLoadDLL("psapi.dll")
	getModuleInformation := psapi.MustFindProc("GetModuleInformation")
	procHandle := windows.CurrentProcess()

	var info moduleInfo
	ret, _, _ := getModuleInformation.Call(uintptr(procHandle), uintptr(handle), uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if ret == 0 {
		err = errors.New("GetModuleInformation failed")
		return
	}

	module = &Module{
		name:   name,
		base:   unsafe.Pointer(info.BaseOfDll),
		size:   uint(info.SizeOfImage),
		handle: handle,
	}
	return
}

// Base get the base pointer of the DLL
func (mod Module) Base() unsafe.Pointer {
	return mod.base
}

type SearchPattern struct {
	Bytes  []uint8
	Ignore []bool
}

type FunctionPattern struct {
	functionName string
	symbolNames  SymbolNameMap
	patterns     PatternMap
	symbolKey    string
	patternKey   string
	addrPointer  unsafe.Pointer
	replaceAddr  unsafe.Pointer
}

func MakeFunctionPattern(functionName string, symbols SymbolNameMap, patterns PatternMap) FunctionPattern {
	return FunctionPattern{functionName: functionName, symbolNames: symbols, patterns: patterns}
}

func (pat *FunctionPattern) Find(module *Module) (foundName string, address unsafe.Pointer, err error) {
	for key, symbolName := range pat.symbolNames {
		var proc uintptr
		proc, err = windows.GetProcAddress(module.handle, symbolName)
		if proc != 0 && err == nil {
			foundName = key
			address = unsafe.Pointer(proc)
			pat.addrPointer = address
			pat.symbolKey = key
			pat.patternKey = ""
			return
		}

		var relAddr uintptr
		relAddr, err = findSym(symbolName)
		if err == nil {
			foundName = key
			address = unsafe.Pointer(uintptr(module.Base()) + relAddr)
			pat.addrPointer = address
			pat.symbolKey = key
			pat.patternKey = ""
			return
		}
	}

	for patternName, pattern := range pat.patterns {
		address, err = findSubstringPattern(pattern, module.base, module.size)
		if err == nil {
			offset := uint(uintptr(address) - uintptr(module.base))
			newSize := module.size - offset
			_, newErr := findSubstringPattern(pattern, unsafe.Pointer(uintptr(address)+1), newSize)
			if newErr == nil {
				// Found a second instance of the pattern, must be an error
				err = ErrNonUniquePattern
			} else {
				foundName = patternName
				pat.addrPointer = address
				pat.symbolKey = ""
				pat.patternKey = patternName
			}
			return
		}
	}
	return
}

func (pat *FunctionPattern) Hook(module *Module, fn unsafe.Pointer) (foundName string, address unsafe.Pointer, err error) {
	foundName, address, err = pat.Find(module)
	if err != nil || foundName == "" {
		return
	}

	if pat.replaceAddr != nil {
		return
	}

	var origFunc uintptr
	if ret := C.MH_CreateHook(C.LPVOID(address), C.LPVOID(fn), (*C.LPVOID)(unsafe.Pointer(&origFunc))); ret != C.MH_OK {
		err = fmt.Errorf("%w: %v", ErrCreateHookFailed, ret)
		return
	}

	if ret := C.MH_EnableHook(C.LPVOID(address)); ret != C.MH_OK {
		err = fmt.Errorf("%w: %v", ErrEnableHookFailed, ret)
		return
	}

	pat.addrPointer = unsafe.Pointer(origFunc)
	pat.replaceAddr = fn
	return
}

func (pat *FunctionPattern) SymbolKey() string {
	return pat.symbolKey
}

func (pat *FunctionPattern) PatternKey() string {
	return pat.patternKey
}

func (pat *FunctionPattern) Name() string {
	return pat.functionName
}

func (pat *FunctionPattern) Address() unsafe.Pointer {
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

func findSubstringPattern(pattern SearchPattern, base unsafe.Pointer, size uint) (addr unsafe.Pointer, err error) {
	patternLen := uint(len(pattern.Bytes))
	for i := uint(0); i < size-patternLen; i++ {
		found := true
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
			return
		}
	}
	err = ErrPatternNotFound
	return
}

var hookInitialized = false

// InitHooks initialise hooks
func InitHooks() (ok bool) {
	if hookInitialized {
		ok = true
		return
	}

	if ret := initSym(); !ret {
		return
	}

	if err := C.MH_Initialize(); err != C.MH_OK {
		return
	}

	hookInitialized = true
	ok = true
	return
}

// CleanupHooks cleanup hooks
func CleanupHooks() {
	if !hookInitialized {
		return
	}

	cleanupSym()

	C.MH_Uninitialize()

	hookInitialized = false
}
