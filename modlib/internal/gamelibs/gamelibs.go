package gamelibs

import (
	"context"
	"hlinspect/internal/engine"
	"hlinspect/internal/hooks"
	"unsafe"
)

import "C"

var Model = &GamelibModel{
	api: &API{
		r: &APIRegistry{},
	},
}

// APIRegistry holds the addresses to game DLL functions.
type APIRegistry struct {
	// HW
	CmdAddCommandWithFlags hooks.FunctionPattern
	CvarRegisterVariable   hooks.FunctionPattern
	MemoryInit             hooks.FunctionPattern
	CmdArgv                hooks.FunctionPattern
	AngleVectors           hooks.FunctionPattern
	PFTracelineDLL         hooks.FunctionPattern

	// CL
	HUDRedraw                   hooks.FunctionPattern
	HUDDrawTransparentTriangles hooks.FunctionPattern
	HUDVidInit                  hooks.FunctionPattern
	HUDReset                    hooks.FunctionPattern

	// Misc
	CCmdHandler unsafe.Pointer
}

// API is a thin interface over the raw game DLL functions. Code that needs to call into
// the game DLLs should do so though this interface. The APIs here should not accept C types,
// nor should they return values in C types.
type API struct {
	r *APIRegistry
}

func (g *API) CmdAddCommand(name string, callback func()) {
	// This implementation is slightly tricky. Here we are taking CCmdHandler set by the gamelib layer.
	// When a command is issued, this is what happens:
	//   hw.dll -> CCmdHandler (C) -> CmdHandler (Go) -> EventHandler.OnCommand
	hooks.CallFuncInts3(g.r.CmdAddCommandWithFlags.Address(), uintptr(unsafe.Pointer(C.CString(name))), uintptr(g.r.CCmdHandler), 2)
}

func (g *API) CmdArgv(arg int) string {
	result := hooks.CallFuncInts1(g.r.CmdArgv.Address(), uintptr(arg))
	return C.GoString((*C.char)(unsafe.Pointer(uintptr(result))))
}

func (g *API) RegisterCVar(cvar *engine.CVar) {
	hooks.CallFuncInts1(g.r.CvarRegisterVariable.Address(), uintptr(cvar.Pointer()))
}

func (g *API) MemoryInit(buf uintptr, size int) {
	hooks.CallFuncInts2(g.r.MemoryInit.Address(), buf, uintptr(size))
}

func (g *API) HUDRedraw(time float32, intermission int) {
	hooks.CallFuncFloatInt(g.r.HUDRedraw.Address(), time, uintptr(intermission))
}

func (g *API) HUDDrawTransparentTriangle() {
	hooks.CallFuncInts0(g.r.HUDDrawTransparentTriangles.Address())
}

func (g *API) HUDVidInit() int {
	return hooks.CallFuncInts0(g.r.HUDVidInit.Address())
}

func (g *API) HUDReset() {
	hooks.CallFuncInts0(g.r.HUDReset.Address())
}

func (g *API) AngleVectors(viewangles [3]float32) (forward, side, up [3]float32) {
	hooks.CallFuncInts4(g.r.AngleVectors.Address(), uintptr(unsafe.Pointer(&viewangles[0])),
		uintptr(unsafe.Pointer(&forward[0])), uintptr(unsafe.Pointer(&side[0])), uintptr(unsafe.Pointer(&up[0])))
	return
}

// TraceLine traces a line and return the hit results
func (g *API) TraceLine(start, end [3]float32, noMonsters int, entToSkip unsafe.Pointer) TraceResult {
	traceResult := TraceResult{}
	hooks.CallFuncInts5(
		g.r.PFTracelineDLL.Address(), uintptr(unsafe.Pointer(&start[0])),
		uintptr(unsafe.Pointer(&end[0])), uintptr(noMonsters),
		uintptr(entToSkip), uintptr(unsafe.Pointer(&traceResult)))
	return traceResult
}

type Handler interface {
	// HW
	MemoryInit(ctx context.Context, buf uintptr, size int)

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
