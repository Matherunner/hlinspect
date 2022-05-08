package cdefs

import (
	"context"
	"hlinspect/internal/logs"
	"unsafe"
)

/*
#include "cdefs.h"
*/
import "C"

var CDefs = struct {
	CCmdHandler                       unsafe.Pointer
	HookedCLCreateMove                unsafe.Pointer
	HookedHUDDrawTransparentTriangles unsafe.Pointer
	HookedHUDRedraw                   unsafe.Pointer
	HookedHUDReset                    unsafe.Pointer
	HookedHUDVidInit                  unsafe.Pointer
	HookedMemoryInit                  unsafe.Pointer
	HookedRClear                      unsafe.Pointer
	HookedRDrawSequentialPoly         unsafe.Pointer
	HookedVFadeAlpha                  unsafe.Pointer
	HookedPMInit                      unsafe.Pointer
	HookedPMPlayerMove                unsafe.Pointer
	CHookedCGraphInitGraph            unsafe.Pointer
	HookedSVExecuteClientMessage      unsafe.Pointer
}{
	CCmdHandler:                       C.CCmdHandler,
	HookedCLCreateMove:                C.HookedCLCreateMove,
	HookedHUDDrawTransparentTriangles: C.HookedHUDDrawTransparentTriangles,
	HookedHUDRedraw:                   C.HookedHUDRedraw,
	HookedHUDReset:                    C.HookedHUDReset,
	HookedHUDVidInit:                  C.HookedHUDVidInit,
	HookedMemoryInit:                  C.HookedMemoryInit,
	HookedRClear:                      C.HookedRClear,
	HookedRDrawSequentialPoly:         C.HookedRDrawSequentialPoly,
	HookedVFadeAlpha:                  C.HookedVFadeAlpha,
	HookedPMInit:                      C.HookedPMInit,
	HookedPMPlayerMove:                C.HookedPMPlayerMove,
	CHookedCGraphInitGraph:            C.CHookedCGraphInitGraph,
	HookedSVExecuteClientMessage:      C.HookedSVExecuteClientMessage,
}

// Handler defines the interface of an event handler that will receive synchronous "events"
// when the hooked functions are called by the game.
//
// This is defined here because the hooked functions define the shape of the handler they want to call.
type Handler interface {
	CLCreateMove(frameTime float32, usercmd unsafe.Pointer, active int)
	HUDDrawTransparentTriangle()
	HUDRedraw(time float32, intermission int)
	HUDReset()
	HUDVidInit() int
	MemoryInit(ctx context.Context, buf uintptr, size int)
	OnCommand()
	RClear()
	RDrawSequentialPoly(surf uintptr, free int)
	VFadeAlpha() int
	PMInit(ppm unsafe.Pointer)
	PMPlayerMove(server int)
	CGraphInitGraph(this unsafe.Pointer)
	SVExecuteClientMessage(cl unsafe.Pointer)
}

var eventHandler Handler

func SetEventHandler(handler Handler) {
	eventHandler = handler
}

// CmdHandler called by C code. This is needed because passing Go function directly to CmdAddCommand doesn't seem to work.
//export CmdHandler
func CmdHandler() {
	defer logs.HandlePanic()
	eventHandler.OnCommand()
}

//export HookedVFadeAlpha
func HookedVFadeAlpha() int {
	defer logs.HandlePanic()
	return eventHandler.VFadeAlpha()
}

//export HookedRClear
func HookedRClear() {
	defer logs.HandlePanic()
	eventHandler.RClear()
}

//export HookedRDrawSequentialPoly
func HookedRDrawSequentialPoly(surf uintptr, free int) {
	defer logs.HandlePanic()
	eventHandler.RDrawSequentialPoly(surf, free)
}

//export HookedMemoryInit
func HookedMemoryInit(buf uintptr, size int) {
	defer logs.HandlePanic()
	eventHandler.MemoryInit(context.TODO(), buf, size)
}

//export HookedCLCreateMove
func HookedCLCreateMove(frameTime float32, usercmd unsafe.Pointer, active int) {
	defer logs.HandlePanic()
	eventHandler.CLCreateMove(frameTime, usercmd, active)
}

//export HookedHUDRedraw
func HookedHUDRedraw(time float32, intermission int32) {
	defer logs.HandlePanic()
	eventHandler.HUDRedraw(time, int(intermission))
}

//export HookedHUDDrawTransparentTriangles
func HookedHUDDrawTransparentTriangles() {
	defer logs.HandlePanic()
	eventHandler.HUDDrawTransparentTriangle()
}

//export HookedHUDVidInit
func HookedHUDVidInit() int {
	defer logs.HandlePanic()
	return eventHandler.HUDVidInit()
}

//export HookedHUDReset
func HookedHUDReset() {
	defer logs.HandlePanic()
	eventHandler.HUDReset()
}

//export HookedPMInit
func HookedPMInit(ppm unsafe.Pointer) {
	defer logs.HandlePanic()
	eventHandler.PMInit(ppm)
}

//export HookedCGraphInitGraph
func HookedCGraphInitGraph(this unsafe.Pointer) {
	defer logs.HandlePanic()
	eventHandler.CGraphInitGraph(this)
}

//export HookedPMPlayerMove
func HookedPMPlayerMove(server int) {
	defer logs.HandlePanic()
	eventHandler.PMPlayerMove(server)
}

//export HookedSVExecuteClientMessage
func HookedSVExecuteClientMessage(cl unsafe.Pointer) {
	defer logs.HandlePanic()
	eventHandler.SVExecuteClientMessage(cl)
}
