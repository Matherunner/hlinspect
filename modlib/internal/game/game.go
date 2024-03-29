package game

import (
	"hlinspect/internal/game/cdefs"
	"hlinspect/internal/game/engine"
	"hlinspect/internal/game/registry"
	"hlinspect/internal/logs"
)

var Model = func() GamelibModel {
	apiReg, err := registry.NewAPI()
	if err != nil {
		logs.DLLLog.Fatalf("unable to initialise API registry: %+v", err)
	}
	return NewGamelibModel(
		NewAPI(apiReg),
		NewGL(),
		NewSync(),
		engine.NewState(),
	)
}()

type Handler = cdefs.Handler

type GamelibModel struct {
	api   *API
	gl    *GL
	sync  *Sync
	state *engine.State
}

func NewGamelibModel(api *API, gl *GL, sync *Sync, state *engine.State) GamelibModel {
	return GamelibModel{
		api:   api,
		gl:    gl,
		sync:  sync,
		state: state,
	}
}

// RegisterEventHandler registers the global handler for all events raised by the game.
// This event handler must be registered before the DLLs are initialised.
func (m *GamelibModel) RegisterEventHandler(handler Handler) {
	cdefs.SetEventHandler(handler)
}

func (m *GamelibModel) InitHWDLL(base string) error {
	return initHWDLL(base)
}

func (m *GamelibModel) InitCLDLL(base string) error {
	return initCLDLL(base)
}

func (m *GamelibModel) InitHLDLL(base string) error {
	return initHLDLL(base)
}

// API returns the gamelib API.
func (m *GamelibModel) API() *API {
	return m.api
}

// GL returns raw OpenGL APIs.
func (m *GamelibModel) GL() *GL {
	return m.gl
}

// Registry returns the API registry.
func (m *GamelibModel) Registry() *registry.API {
	return m.api.Registry()
}

func (m *GamelibModel) Sync() *Sync {
	return m.sync
}

// S returns the engine state data.
func (m *GamelibModel) S() *engine.State {
	return m.state
}
