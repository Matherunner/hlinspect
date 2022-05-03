package main

import (
	"context"
	"hlinspect/internal/events"
	"hlinspect/internal/feed"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/hlrpc"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"io"
	"net"
	"path/filepath"
	"unsafe"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"golang.org/x/sys/windows"
)

/*
#include "dllmain.h"
*/
import "C"

var kernelDLL *hooks.Module

var libraryInitializers = map[string]func(base string) error{
	"hl.dll":     gamelibs.Model.InitHLDLL,
	"opfor.dll":  gamelibs.Model.InitHLDLL,
	"cz.dll":     gamelibs.Model.InitHLDLL,
	"gunman.dll": gamelibs.Model.InitHLDLL,
	"wanted.dll": gamelibs.Model.InitHLDLL,
	"hw.dll":     gamelibs.Model.InitHWDLL,
	"client.dll": gamelibs.Model.InitCLDLL,
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

	initLoadLibraryHooks()

	gamelibs.Model.RegisterEventHandler(events.NewHandler())

	for base, initializer := range libraryInitializers {
		logs.DLLLog.Debugf("Initialising %v", base)
		if err := initializer(base); err != nil {
			logs.DLLLog.Warningf("Unable to initialise %v: %v", base, err)
		}
	}

	go func() {
		ln, err := net.Listen("tcp4", "0.0.0.0:32002")
		if err != nil {
			panic(err)
		}

		for {
			logs.DLLLog.Debugf("waiting to accept Conn: %+v", ln.Addr())
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}

			logs.DLLLog.Debugf("accepted a new connection %+v", conn)

			serveHLRPC(conn)
		}
	}()

	go feed.Serve()
}

type halflifeServer struct {
}

func (s *halflifeServer) GetFullPlayerState(ctx context.Context, call hlrpc.HalfLife_getFullPlayerState) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}

	fullState, err := hlrpc.NewRootFullPlayerState(seg)
	if err != nil {
		panic(err)
	}
	fullState.SetVelocityX(320)
	fullState.SetVelocityY(640)
	fullState.SetWaterLevel(2)
	fullState.SetDuckState(hlrpc.DuckState_ducked)

	res.SetState(fullState)
	return nil
}

func serveHLRPC(rwc io.ReadWriteCloser) {
	defer rwc.Close()

	main := hlrpc.HalfLife_ServerToClient(&halflifeServer{}, nil)

	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		BootstrapClient: main.Client,
	})
	defer conn.Close()

	select {
	case <-conn.Done():
		return
	}
}

// OnProcessDetach called from DllMain on process detach
//export OnProcessDetach
func OnProcessDetach() {
	hooks.CleanupHooks()
}

func main() {}
