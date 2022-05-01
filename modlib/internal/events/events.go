package events

import (
	"context"
	"hlinspect/internal/cmd"
	"hlinspect/internal/cvar"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gamelibs/graphics"
	"hlinspect/internal/gamelibs/hw"
	"hlinspect/internal/gl"
)

type handler struct {
}

func NewHandler() gamelibs.Handler {
	return &handler{}
}

func (h *handler) MemoryInit(ctx context.Context, buf uintptr, size int) {
	gamelibs.Model.API().MemoryInit(buf, size)

	gamelibs.Model.API().RegisterCVar(&cvar.Wallhack)
	gamelibs.Model.API().RegisterCVar(&cvar.Wallhack)
	gamelibs.Model.API().RegisterCVar(&cvar.WallhackAlpha)
	gamelibs.Model.API().RegisterCVar(&cvar.FadeRemove)
	gamelibs.Model.API().RegisterCVar(&cvar.Nodes)
	gamelibs.Model.API().RegisterCVar(&cvar.NodeLinks)
	gamelibs.Model.API().RegisterCVar(&cvar.Cine)
	gamelibs.Model.API().RegisterCVar(&cvar.CinePossess)

	for name, handler := range cmd.CommandHandlerByName {
		gamelibs.Model.API().CmdAddCommand(name, handler)
	}
}

func (h *handler) OnCommand() {
	name := gamelibs.Model.API().CmdArgv(0)
	if handler, ok := cmd.CommandHandlerByName[name]; ok {
		handler()
	}
}

func (h *handler) HUDRedraw(time float32, intermission int) {
	gamelibs.Model.API().HUDRedraw(time, intermission)

	graphics.DrawHUD(time, intermission)
}

func (h *handler) HUDDrawTransparentTriangle() {
	gamelibs.Model.API().HUDDrawTransparentTriangle()

	gl.Disable(gl.Texture2D)
	graphics.DrawTriangles()
	gl.Enable(gl.Texture2D)
	hw.TriGLRenderMode(hw.KRenderNormal)
}

func (h *handler) HUDVidInit() int {
	ret := gamelibs.Model.API().HUDVidInit()

	screenInfo := hw.GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)

	return ret
}

func (h *handler) HUDReset() {
	gamelibs.Model.API().HUDReset()

	screenInfo := hw.GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)
}
