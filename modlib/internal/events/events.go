package events

import (
	"context"
	"hlinspect/internal/cmd"
	"hlinspect/internal/cvar"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gamelibs/graphics"
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

	for name := range cmd.CommandHandlerByName {
		gamelibs.Model.API().CmdAddCommand(name)
	}
}

func (h *handler) OnCommand() {
	name := gamelibs.Model.API().CmdArgv(0)
	if handler, ok := cmd.CommandHandlerByName[name]; ok {
		handler()
	}
}

func (h *handler) VFadeAlpha() int {
	if cvar.FadeRemove.Float32() != 0 {
		return 0
	}
	return gamelibs.Model.API().VFadeAlpha()
}

func (h *handler) RClear() {
	if cvar.Wallhack.Float32() != 0 {
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.ColorBufferBit)
	}
	gamelibs.Model.API().RClear()
}

func (h *handler) RDrawSequentialPoly(surf uintptr, free int) {
	if cvar.Wallhack.Float32() == 0 {
		gamelibs.Model.API().RDrawSequentialPoly(surf, free)
		return
	}

	gl.Enable(gl.Blend)
	gl.DepthMask(false)
	gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
	gl.Color4f(1, 1, 1, cvar.WallhackAlpha.Float32())

	gamelibs.Model.API().RDrawSequentialPoly(surf, free)

	gl.DepthMask(true)
	gl.Disable(gl.Blend)
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
	gamelibs.Model.API().TriGLRenderMode(gamelibs.KRenderNormal)
}

func (h *handler) HUDVidInit() int {
	ret := gamelibs.Model.API().HUDVidInit()

	screenInfo := gamelibs.Model.API().GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)

	return ret
}

func (h *handler) HUDReset() {
	gamelibs.Model.API().HUDReset()

	screenInfo := gamelibs.Model.API().GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)
}
