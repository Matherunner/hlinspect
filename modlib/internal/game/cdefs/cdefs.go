package cdefs

import (
	"context"
	"unsafe"
)

/*
#include "cdefs.h"
*/
import "C"

var CDefs = struct {
	CCmdHandler                       unsafe.Pointer
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
}{
	CCmdHandler:                       C.CCmdHandler,
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
}

// Handler defines the interface of an event handler that will receive synchronous "events"
// when the hooked functions are called by the game.
//
// This is defined here because the hooked functions define the shape of the handler they want to call.
type Handler interface {
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
}

var eventHandler Handler

func SetEventHandler(handler Handler) {
	eventHandler = handler
}

// CmdHandler called by C code. This is needed because passing Go function directly to CmdAddCommand doesn't seem to work.
//export CmdHandler
func CmdHandler() {
	eventHandler.OnCommand()
}

//export HookedVFadeAlpha
func HookedVFadeAlpha() int {
	return eventHandler.VFadeAlpha()
}

//export HookedRClear
func HookedRClear() {
	eventHandler.RClear()
}

//export HookedRDrawSequentialPoly
func HookedRDrawSequentialPoly(surf uintptr, free int) {
	eventHandler.RDrawSequentialPoly(surf, free)
}

//export HookedMemoryInit
func HookedMemoryInit(buf uintptr, size int) {
	eventHandler.MemoryInit(context.TODO(), buf, size)
}

//export HookedHUDRedraw
func HookedHUDRedraw(time float32, intermission int32) {
	eventHandler.HUDRedraw(time, int(intermission))
}

//export HookedHUDDrawTransparentTriangles
func HookedHUDDrawTransparentTriangles() {
	eventHandler.HUDDrawTransparentTriangle()
}

//export HookedHUDVidInit
func HookedHUDVidInit() int {
	return eventHandler.HUDVidInit()
}

//export HookedHUDReset
func HookedHUDReset() {
	eventHandler.HUDReset()
}

//export HookedPMInit
func HookedPMInit(ppm unsafe.Pointer) {
	eventHandler.PMInit(ppm)
}

//export HookedCGraphInitGraph
func HookedCGraphInitGraph(this unsafe.Pointer) {
	eventHandler.CGraphInitGraph(this)
}

//export HookedPMPlayerMove
func HookedPMPlayerMove(server int) {
	eventHandler.PMPlayerMove(server)
}
