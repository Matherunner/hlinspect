package game

import (
	"hlinspect/internal/game/cdefs"
	"hlinspect/internal/game/engine"
	"hlinspect/internal/game/registry"
	"hlinspect/internal/hooks"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

// API is a thin interface over the raw game DLL functions. Code that needs to call into
// the game DLLs should do so though this interface. The APIs here should not accept C types,
// nor should they return values in C types.
type API struct {
	r *registry.API
}

func NewAPI(apiRegistry *registry.API) *API {
	return &API{
		r: apiRegistry,
	}
}

func (api *API) Registry() *registry.API {
	return api.r
}

func (api *API) BuildNumber() int {
	return hooks.CallFuncInts0(api.r.BuildNumber.Ptr())
}

func (api *API) GetScreenInfo() ScreenInfo {
	screenInfo := ScreenInfo{}
	screenInfo.Size = int32(unsafe.Sizeof(screenInfo))
	hooks.CallFuncInts1(api.r.HudGetScreenInfo.Ptr(), uintptr(unsafe.Pointer(&screenInfo)))
	return screenInfo
}

func (api *API) DrawString(x, y int, text string) {
	ctext := unsafe.Pointer(C.CString(text))
	defer C.free(ctext)
	hooks.CallFuncInts3(api.r.DrawString.Ptr(), uintptr(x), uintptr(y), uintptr(ctext))
}

func (api *API) VGUI2DrawSetTextColorAlpha(r, g, b, a int) {
	hooks.CallFuncInts4(api.r.VGUI2DrawSetTextColorAlpha.Ptr(), uintptr(r), uintptr(g), uintptr(b), uintptr(a))
}

func (api *API) VFadeAlpha() int {
	return hooks.CallFuncInts0(api.r.VFadeAlpha.Ptr())
}

func (api *API) RClear() {
	hooks.CallFuncInts0(api.r.RClear.Ptr())
}

func (api *API) RDrawSequentialPoly(surf uintptr, free int) {
	hooks.CallFuncInts2(api.r.RDrawSequentialPoly.Ptr(), surf, uintptr(free))
}

func (api *API) TriGLRenderMode(mode int) {
	hooks.CallFuncInts1(api.r.TriGLRenderMode.Ptr(), uintptr(mode))
}

func (api *API) TriGLBegin(primitive int) {
	hooks.CallFuncInts1(api.r.TriGLBegin.Ptr(), uintptr(primitive))
}

func (api *API) TriGLEnd() {
	hooks.CallFuncInts0(api.r.TriGLEnd.Ptr())
}

func (api *API) TriGLColor4f(r, g, b, a float32) {
	hooks.CallFuncFloats4(api.r.TriGLColor4f.Ptr(), r, g, b, a)
}

func (api *API) TriGLVertex3fv(v [3]float32) {
	hooks.CallFuncInts1(api.r.TriGLVertex3fv.Ptr(), uintptr(unsafe.Pointer(&v[0])))
}

func (api *API) TriGLCullFace(style int) {
	hooks.CallFuncInts1(api.r.TriGLCullFace.Ptr(), uintptr(style))
}

// ScreenTransform ScreenTransform, similar to WorldToScreen in TriAPI
func (api *API) ScreenTransform(point [3]float32) (screen [3]float32, clipped bool) {
	clipped = hooks.CallFuncInts2(api.r.ScreenTransform.Ptr(), uintptr(unsafe.Pointer(&point[0])), uintptr(unsafe.Pointer(&screen[0]))) != 0
	return
}

func (api *API) PFCheckClientI(edict unsafe.Pointer) uintptr {
	return uintptr(hooks.CallFuncInts1(api.r.PFCheckClientI.Ptr(), uintptr(edict)))
}

// CmdAddCommand registers the given name to a common command handler.
func (api *API) CmdAddCommand(name string) {
	// This implementation is slightly tricky because we can't register a Go function.
	// We have to register the same CCmdHandler set by the gamelib layer for every command.
	// When a command is issued, this is what happens:
	//   hw.dll -> CCmdHandler (C) -> CmdHandler (Go) -> EventHandler.OnCommand
	// Then the implementation of OnCommand should distinguish which command is actually called using Cmd_Argv(0).
	// The name does not need to be freed because the registered command is global.
	hooks.CallFuncInts3(api.r.CmdAddCommandWithFlags.Ptr(), uintptr(unsafe.Pointer(C.CString(name))), uintptr(cdefs.CDefs.CCmdHandler), 2)
}

func (api *API) CmdArgv(arg int) string {
	result := hooks.CallFuncInts1RetPtr(api.r.CmdArgv.Ptr(), uintptr(arg))
	return C.GoString((*C.char)(result))
}

func (api *API) RegisterCVar(cvar *engine.CVar) {
	hooks.CallFuncInts1(api.r.CvarRegisterVariable.Ptr(), uintptr(cvar.Pointer()))
}

func (api *API) MemoryInit(buf uintptr, size int) {
	hooks.CallFuncInts2(api.r.MemoryInit.Ptr(), buf, uintptr(size))
}

func (api *API) HUDRedraw(time float32, intermission int) {
	hooks.CallFuncFloatInt(api.r.HUDRedraw.Ptr(), time, uintptr(intermission))
}

func (api *API) HUDDrawTransparentTriangle() {
	hooks.CallFuncInts0(api.r.HUDDrawTransparentTriangles.Ptr())
}

func (api *API) HUDVidInit() int {
	return hooks.CallFuncInts0(api.r.HUDVidInit.Ptr())
}

func (api *API) HUDReset() {
	hooks.CallFuncInts0(api.r.HUDReset.Ptr())
}

func (api *API) AngleVectors(viewangles [3]float32) (forward, side, up [3]float32) {
	hooks.CallFuncInts4(api.r.AngleVectors.Ptr(), uintptr(unsafe.Pointer(&viewangles[0])),
		uintptr(unsafe.Pointer(&forward[0])), uintptr(unsafe.Pointer(&side[0])), uintptr(unsafe.Pointer(&up[0])))
	return
}

func (api *API) TraceLine(start, end [3]float32, noMonsters int, entToSkip unsafe.Pointer) (result TraceResult) {
	hooks.CallFuncInts5(
		api.r.PFTracelineDLL.Ptr(), uintptr(unsafe.Pointer(&start[0])),
		uintptr(unsafe.Pointer(&end[0])), uintptr(noMonsters),
		uintptr(entToSkip), uintptr(unsafe.Pointer(&result)))
	return
}

func (api *API) CSoundEntActiveList() int32 {
	return int32(hooks.CallFuncInts0(api.r.CSoundEntActiveList.Ptr()))
}

func (api *API) CSoundEntSoundPointerForIndex(index int32) unsafe.Pointer {
	return hooks.CallFuncInts1RetPtr(api.r.CSoundEntSoundPointerForIndex.Ptr(), uintptr(index))
}

func (api *API) CBaseMonsterPBestSound(this unsafe.Pointer) unsafe.Pointer {
	return hooks.CallFuncThisInts0RetPtr(api.r.CBaseMonsterPBestSound.Ptr(), uintptr(this))
}

func (api *API) CGraphInitGraph(this unsafe.Pointer) {
	hooks.CallFuncThisInts0(api.r.CGraphInitGraph.Ptr(), uintptr(this))
}

func (api *API) PMInit(ppm unsafe.Pointer) {
	hooks.CallFuncInts1(api.r.PMInit.Ptr(), uintptr(ppm))
}

func (api *API) PMPlayerMove(server int) {
	hooks.CallFuncInts1(api.r.PMPlayerMove.Ptr(), uintptr(server))
}

func (api *API) CLCreateMove(frameTime float32, usercmd unsafe.Pointer, active int) {
	hooks.CallFuncFloatInts2(api.r.CLCreateMove.Ptr(), frameTime, uintptr(usercmd), uintptr(active))
}

func (api *API) CbufInsertText(text string) {
	cs := unsafe.Pointer(C.CString(text))
	defer C.free(cs)
	hooks.CallFuncInts1(api.r.CbufInsertText.Ptr(), uintptr(cs))
}

func (api *API) WriteDestParm(dest int) unsafe.Pointer {
	return hooks.CallFuncInts1RetPtr(api.r.WriteDestParm.Ptr(), uintptr(dest))
}

func (api *API) SVExecuteClientMessage(cl unsafe.Pointer) {
	hooks.CallFuncInts1(api.r.SVExecuteClientMessage.Ptr(), uintptr(cl))
}
