package gamelibs

import (
	"context"
	"hlinspect/internal/engine"
	"hlinspect/internal/hooks"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

var Model = &GamelibModel{
	api: &API{
		r: &APIRegistry{},
	},
}

// APIRegistry holds the addresses to game DLL functions.
type APIRegistry struct {
	// HW
	AngleVectors               hooks.FunctionPattern
	BuildNumber                hooks.FunctionPattern
	CmdAddCommandWithFlags     hooks.FunctionPattern
	CmdArgv                    hooks.FunctionPattern
	CvarRegisterVariable       hooks.FunctionPattern
	DrawString                 hooks.FunctionPattern
	HostAutoSaveF              hooks.FunctionPattern
	HostNoclipF                hooks.FunctionPattern
	HudGetScreenInfo           hooks.FunctionPattern
	MemoryInit                 hooks.FunctionPattern
	PFCheckClientI             hooks.FunctionPattern
	PFTracelineDLL             hooks.FunctionPattern
	RClear                     hooks.FunctionPattern
	RDrawSequentialPoly        hooks.FunctionPattern
	ScreenTransform            hooks.FunctionPattern
	TriGLBegin                 hooks.FunctionPattern
	TriGLColor4f               hooks.FunctionPattern
	TriGLCullFace              hooks.FunctionPattern
	TriGLEnd                   hooks.FunctionPattern
	TriGLRenderMode            hooks.FunctionPattern
	TriGLVertex3fv             hooks.FunctionPattern
	VFadeAlpha                 hooks.FunctionPattern
	VGUI2DrawSetTextColorAlpha hooks.FunctionPattern
	WorldTransform             hooks.FunctionPattern

	// CL
	HUDRedraw                   hooks.FunctionPattern
	HUDDrawTransparentTriangles hooks.FunctionPattern
	HUDVidInit                  hooks.FunctionPattern
	HUDReset                    hooks.FunctionPattern

	// HL

	// Misc
	CCmdHandler unsafe.Pointer
}

// API is a thin interface over the raw game DLL functions. Code that needs to call into
// the game DLLs should do so though this interface. The APIs here should not accept C types,
// nor should they return values in C types.
type API struct {
	r *APIRegistry
}

func (api *API) BuildNumber() int {
	return hooks.CallFuncInts0(api.r.BuildNumber.Address())
}

func (api *API) GetScreenInfo() ScreenInfo {
	screenInfo := ScreenInfo{}
	screenInfo.Size = int32(unsafe.Sizeof(screenInfo))
	hooks.CallFuncInts1(api.r.HudGetScreenInfo.Address(), uintptr(unsafe.Pointer(&screenInfo)))
	return screenInfo
}

func (api *API) DrawString(x, y int, text string) {
	ctext := unsafe.Pointer(C.CString(text))
	defer C.free(ctext)
	hooks.CallFuncInts3(api.r.DrawString.Address(), uintptr(x), uintptr(y), uintptr(ctext))
}

func (api *API) VGUI2DrawSetTextColorAlpha(r, g, b, a int) {
	hooks.CallFuncInts4(api.r.VGUI2DrawSetTextColorAlpha.Address(), uintptr(r), uintptr(g), uintptr(b), uintptr(a))
}

func (api *API) VFadeAlpha() int {
	return hooks.CallFuncInts0(api.r.VFadeAlpha.Address())
}

func (api *API) RClear() {
	hooks.CallFuncInts0(api.r.RClear.Address())
}

func (api *API) RDrawSequentialPoly(surf uintptr, free int) {
	hooks.CallFuncInts2(api.r.RDrawSequentialPoly.Address(), surf, uintptr(free))
}

func (api *API) TriGLRenderMode(mode int) {
	hooks.CallFuncInts1(api.r.TriGLRenderMode.Address(), uintptr(mode))
}

func (api *API) TriGLBegin(primitive int) {
	hooks.CallFuncInts1(api.r.TriGLBegin.Address(), uintptr(primitive))
}

func (api *API) TriGLEnd() {
	hooks.CallFuncInts0(api.r.TriGLEnd.Address())
}

func (api *API) TriGLColor4f(r, g, b, a float32) {
	hooks.CallFuncFloats4(api.r.TriGLColor4f.Address(), r, g, b, a)
}

func (api *API) TriGLVertex3fv(v [3]float32) {
	hooks.CallFuncInts1(api.r.TriGLVertex3fv.Address(), uintptr(unsafe.Pointer(&v[0])))
}

func (api *API) TriGLCullFace(style int) {
	hooks.CallFuncInts1(api.r.TriGLCullFace.Address(), uintptr(style))
}

// ScreenTransform ScreenTransform, similar to WorldToScreen in TriAPI
func (api *API) ScreenTransform(point [3]float32) (screen [3]float32, clipped bool) {
	clipped = hooks.CallFuncInts2(api.r.ScreenTransform.Address(), uintptr(unsafe.Pointer(&point[0])), uintptr(unsafe.Pointer(&screen[0]))) != 0
	return
}

func (api *API) PFCheckClientI(edict unsafe.Pointer) uintptr {
	return uintptr(hooks.CallFuncInts1(api.r.PFCheckClientI.Address(), uintptr(edict)))
}

// CmdAddCommand registers the given name to a common command handler.
func (api *API) CmdAddCommand(name string) {
	// This implementation is slightly tricky because we can't register a Go function.
	// We have to register the same CCmdHandler set by the gamelib layer for every command.
	// When a command is issued, this is what happens:
	//   hw.dll -> CCmdHandler (C) -> CmdHandler (Go) -> EventHandler.OnCommand
	// Then the implementation of OnCommand should distinguish which command is actually called using Cmd_Argv(0).
	// The name does not need to be freed because the registered command is global.
	hooks.CallFuncInts3(api.r.CmdAddCommandWithFlags.Address(), uintptr(unsafe.Pointer(C.CString(name))), uintptr(api.r.CCmdHandler), 2)
}

func (api *API) CmdArgv(arg int) string {
	result := hooks.CallFuncInts1(api.r.CmdArgv.Address(), uintptr(arg))
	return C.GoString((*C.char)(unsafe.Pointer(uintptr(result))))
}

func (api *API) RegisterCVar(cvar *engine.CVar) {
	hooks.CallFuncInts1(api.r.CvarRegisterVariable.Address(), uintptr(cvar.Pointer()))
}

func (api *API) MemoryInit(buf uintptr, size int) {
	hooks.CallFuncInts2(api.r.MemoryInit.Address(), buf, uintptr(size))
}

func (api *API) HUDRedraw(time float32, intermission int) {
	hooks.CallFuncFloatInt(api.r.HUDRedraw.Address(), time, uintptr(intermission))
}

func (api *API) HUDDrawTransparentTriangle() {
	hooks.CallFuncInts0(api.r.HUDDrawTransparentTriangles.Address())
}

func (api *API) HUDVidInit() int {
	return hooks.CallFuncInts0(api.r.HUDVidInit.Address())
}

func (api *API) HUDReset() {
	hooks.CallFuncInts0(api.r.HUDReset.Address())
}

func (api *API) AngleVectors(viewangles [3]float32) (forward, side, up [3]float32) {
	hooks.CallFuncInts4(api.r.AngleVectors.Address(), uintptr(unsafe.Pointer(&viewangles[0])),
		uintptr(unsafe.Pointer(&forward[0])), uintptr(unsafe.Pointer(&side[0])), uintptr(unsafe.Pointer(&up[0])))
	return
}

// TraceLine traces a line and return the hit results
func (api *API) TraceLine(start, end [3]float32, noMonsters int, entToSkip unsafe.Pointer) TraceResult {
	traceResult := TraceResult{}
	hooks.CallFuncInts5(
		api.r.PFTracelineDLL.Address(), uintptr(unsafe.Pointer(&start[0])),
		uintptr(unsafe.Pointer(&end[0])), uintptr(noMonsters),
		uintptr(entToSkip), uintptr(unsafe.Pointer(&traceResult)))
	return traceResult
}

type Handler interface {
	// HW
	MemoryInit(ctx context.Context, buf uintptr, size int)
	VFadeAlpha() int
	RClear()
	RDrawSequentialPoly(surf uintptr, free int)

	// CL
	HUDRedraw(time float32, intermission int)
	HUDDrawTransparentTriangle()
	HUDVidInit() int
	HUDReset()

	// Misc
	OnCommand()
}

type GamelibModel struct {
	eventHandler Handler
	api          *API
}

// RegisterEventHandler registers the global handler for all events raised by the gamelibs.
func (m *GamelibModel) RegisterEventHandler(handler Handler) {
	m.eventHandler = handler
}

// EventHandler returns the event handler registered using RegisterEventHandler.
func (m *GamelibModel) EventHandler() Handler {
	return m.eventHandler
}

// API returns the gamelib API.
func (m *GamelibModel) API() *API {
	return m.api
}

// Registry returns the API registry.
func (m *GamelibModel) Registry() *APIRegistry {
	return m.api.r
}
