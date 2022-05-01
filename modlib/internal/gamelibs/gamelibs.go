package gamelibs

import (
	"hlinspect/internal/gamelibs/cdefs"
)

var Model = &GamelibModel{
	api: &API{
		r: &APIRegistry{},
	},
}

type Handler = cdefs.Handler

type GamelibModel struct {
	api *API
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

// Registry returns the API registry.
func (m *GamelibModel) Registry() *APIRegistry {
	return m.api.r
}
