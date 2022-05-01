package gamelibs

import (
	"context"
	"hlinspect/internal/engine"
	"hlinspect/internal/hooks"
)

var Model = &GamelibModel{
	api: &API{
		r: &APIRegistry{},
	},
}

// APIRegistry holds the addresses to game DLL functions.
type APIRegistry struct {
	// HW
	CvarRegisterVariable hooks.FunctionPattern
	MemoryInit           hooks.FunctionPattern

	// CL
	HUDRedrawPattern                   hooks.FunctionPattern
	HUDDrawTransparentTrianglesPattern hooks.FunctionPattern
	HUDVidInitPattern                  hooks.FunctionPattern
	HUDResetPattern                    hooks.FunctionPattern
}

// API is a thin interface over the raw game DLL functions. Code that needs to call into
// the game DLLs should do so though this interface. The APIs here should not accept C types,
// nor should they return values in C types.
type API struct {
	r *APIRegistry
}

func (g *API) RegisterCVar(cvar *engine.CVar) {
	hooks.CallFuncInts1(g.r.CvarRegisterVariable.Address(), uintptr(cvar.Pointer()))
}

func (g *API) MemoryInit(buf uintptr, size int) {
	hooks.CallFuncInts2(g.r.MemoryInit.Address(), buf, uintptr(size))
}

func (g *API) HUDRedraw(time float32, intermission int) {
	hooks.CallFuncFloatInt(g.r.HUDRedrawPattern.Address(), time, uintptr(intermission))
}

func (g *API) HUDDrawTransparentTriangle() {
	hooks.CallFuncInts0(g.r.HUDDrawTransparentTrianglesPattern.Address())
}

func (g *API) HUDVidInit() int {
	return hooks.CallFuncInts0(g.r.HUDVidInitPattern.Address())
}

func (g *API) HUDReset() {
	hooks.CallFuncInts0(g.r.HUDResetPattern.Address())
}

type Handler interface {
	// HW
	MemoryInit(ctx context.Context, buf uintptr, size int)

	// CL
	HUDRedraw(time float32, intermission int)
	HUDDrawTransparentTriangle()
	HUDVidInit() int
	HUDReset()
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
