package main

import (
	"hlinspect/internal/feed"
	"hlinspect/internal/game"
	"hlinspect/internal/handlers"
	"hlinspect/internal/hlrpc"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
#include "dllmain.h"
*/
import "C"

var kernelDLL *hooks.Module

var libraryInitializers = map[string]func(base string) error{
	"hl.dll":     game.Model.InitHLDLL,
	"opfor.dll":  game.Model.InitHLDLL,
	"cz.dll":     game.Model.InitHLDLL,
	"gunman.dll": game.Model.InitHLDLL,
	"wanted.dll": game.Model.InitHLDLL,
	"hw.dll":     game.Model.InitHWDLL,
	"client.dll": game.Model.InitCLDLL,
}

var loadLibraryAPattern = hooks.NewFunctionPattern("LoadLibraryA", hooks.SymbolNameMap{"Windows": "LoadLibraryA"}, nil)
var loadLibraryWPattern = hooks.NewFunctionPattern("LoadLibraryW", hooks.SymbolNameMap{"Windows": "LoadLibraryW"}, nil)

// GetLoadLibraryAAddr called by C to get the address of the original LoadLibraryA
//export GetLoadLibraryAAddr
func GetLoadLibraryAAddr() uintptr {
	return uintptr(loadLibraryAPattern.Ptr())
}

// LoadLibraryACallback called by C when the library has been loaded successfully
//export LoadLibraryACallback
func LoadLibraryACallback(fileName C.LPCSTR) {
	onLibraryLoaded(C.GoString(fileName))
}

// GetLoadLibraryWAddr called by C to get the address of the original LoadLibraryW
//export GetLoadLibraryWAddr
func GetLoadLibraryWAddr() uintptr {
	return uintptr(loadLibraryWPattern.Ptr())
}

// LoadLibraryWCallback called by C when the library has been loaded successfully
//export LoadLibraryWCallback
func LoadLibraryWCallback(fileName C.LPCWSTR) {
	onLibraryLoaded(windows.UTF16PtrToString((*uint16)(unsafe.Pointer(fileName))))
}

func onLibraryLoaded(fileName string) {
	hooks.RefreshModuleList()
	base := filepath.Base(fileName)
	if initializer, ok := libraryInitializers[base]; ok {
		if err := initializer(base); err != nil {
			logs.DLLLog.Warningf("Unable to hook %v when loaded: %v", base, err)
		} else {
			logs.DLLLog.Debugf("Initialised %v", base)
		}
	}
}

func initLoadLibraryHooks() {
	var err error
	kernelDLL, err = hooks.NewModule("kernel32.dll")
	if err != nil {
		logs.DLLLog.Panic("Unable to initialise kernel32.dll")
	}

	logs.DLLLog.Debug("Hooking LoadLibraryA")
	_, _, err = loadLibraryAPattern.Hook(kernelDLL, C.HookedLoadLibraryA)
	if err != nil {
		logs.DLLLog.Panicf("Unable to hook LoadLibraryA: %v", err)
	}

	logs.DLLLog.Debug("Hooking LoadLibraryW")
	_, _, err = loadLibraryWPattern.Hook(kernelDLL, C.HookedLoadLibraryW)
	if err != nil {
		logs.DLLLog.Panicf("Unable to hook LoadLibraryW: %v", err)
	}
}

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

	game.Model.RegisterEventHandler(handlers.NewGameHandler())

	initLoadLibraryHooks()

	for base, initializer := range libraryInitializers {
		logs.DLLLog.Debugf("Initialising %v", base)
		if err := initializer(base); err != nil {
			logs.DLLLog.Warningf("Unable to initialise %v: %v", base, err)
		}
	}

	go func() {
		err := hlrpc.Serve(handlers.NewHLRPCHandler())
		logs.DLLLog.Errorf("hlrpc failed: %+v", err)
	}()

	go feed.Serve()
}

// OnProcessDetach called from DllMain on process detach
//export OnProcessDetach
func OnProcessDetach() {
	hooks.CleanupHooks()
}

func main() {}
