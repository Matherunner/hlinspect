package gamelibs

import (
	"hlinspect/internal/gamelibs/cdefs"
	"hlinspect/internal/gamelibs/registry"
	"hlinspect/internal/logs"
)

var Model = func() GamelibModel {
	apiReg, err := registry.NewAPI()
	if err != nil {
		logs.DLLLog.Fatalf("unable to initialise API registry: %+v", err)
	}
	return NewGamelibModel(NewAPI(apiReg), NewGL())
}()

type Handler = cdefs.Handler

type GamelibModel struct {
	api *API
	gl  *GL
}

func NewGamelibModel(api *API, gl *GL) GamelibModel {
	return GamelibModel{
		api: api,
		gl:  gl,
	}
}

// RegisterEventHandler registers the global handler for all events raised by the gamelibs.
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
