package handlers

import (
	"context"
	"hlinspect/internal/cmd"
	"hlinspect/internal/cvar"
	"hlinspect/internal/engine"
	"hlinspect/internal/feed"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/graphics"
	"hlinspect/internal/logs"
	"hlinspect/internal/proto"
	"unsafe"
)

type gameHandler struct {
}

func NewGameHandler() *gameHandler {
	return &gameHandler{}
}

func (h *gameHandler) MemoryInit(ctx context.Context, buf uintptr, size int) {
	gamelibs.Model.API().MemoryInit(buf, size)

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

func (h *gameHandler) OnCommand() {
	name := gamelibs.Model.API().CmdArgv(0)
	if handler, ok := cmd.CommandHandlerByName[name]; ok {
		handler()
	}
}

func (h *gameHandler) VFadeAlpha() int {
	if cvar.FadeRemove.Float32() != 0 {
		return 0
	}
	return gamelibs.Model.API().VFadeAlpha()
}

func (h *gameHandler) RClear() {
	if cvar.Wallhack.Float32() != 0 {
		gamelibs.Model.GL().ClearColor(0, 0, 0, 1)
		gamelibs.Model.GL().Clear(gamelibs.GLColorBufferBit)
	}
	gamelibs.Model.API().RClear()
}

func (h *gameHandler) RDrawSequentialPoly(surf uintptr, free int) {
	if cvar.Wallhack.Float32() == 0 {
		gamelibs.Model.API().RDrawSequentialPoly(surf, free)
		return
	}

	gamelibs.Model.GL().Enable(gamelibs.GLBlend)
	gamelibs.Model.GL().DepthMask(false)
	gamelibs.Model.GL().BlendFunc(gamelibs.GLSrcAlpha, gamelibs.GLOneMinusSrcAlpha)
	gamelibs.Model.GL().Color4f(1, 1, 1, cvar.WallhackAlpha.Float32())

	gamelibs.Model.API().RDrawSequentialPoly(surf, free)

	gamelibs.Model.GL().DepthMask(true)
	gamelibs.Model.GL().Disable(gamelibs.GLBlend)
}

func (h *gameHandler) HUDRedraw(time float32, intermission int) {
	gamelibs.Model.API().HUDRedraw(time, intermission)

	graphics.DrawHUD(time, intermission)
}

func (h *gameHandler) HUDDrawTransparentTriangle() {
	gamelibs.Model.API().HUDDrawTransparentTriangle()

	gamelibs.Model.GL().Disable(gamelibs.GLTexture2D)
	graphics.DrawTriangles()
	gamelibs.Model.GL().Enable(gamelibs.GLTexture2D)
	gamelibs.Model.API().TriGLRenderMode(gamelibs.KRenderNormal)
}

func (h *gameHandler) HUDVidInit() int {
	ret := gamelibs.Model.API().HUDVidInit()

	screenInfo := gamelibs.Model.API().GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)

	return ret
}

func (h *gameHandler) HUDReset() {
	gamelibs.Model.API().HUDReset()

	screenInfo := gamelibs.Model.API().GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)
}

func (h *gameHandler) PMInit(ppm unsafe.Pointer) {
	gamelibs.Model.API().PMInit(ppm)

	engine.Engine.SetPPMove(ppm)
	logs.DLLLog.Debugf("Set PPMOVE with address = %x", ppm)
}

func (h *gameHandler) PMPlayerMove(server int) {
	binary, err := proto.Serialize(&proto.PMove{
		Stage:        proto.PMoveStagePre,
		Velocity:     engine.Engine.PMoveVelocity(),
		Position:     engine.Engine.PMovePosition(),
		Viewangles:   engine.Engine.PMoveViewangles(),
		Basevelocity: engine.Engine.PMoveBasevelocity(),
		FSU:          engine.Engine.PMoveCmdFSU(),
		Punchangles:  engine.Engine.PMovePunchangles(),
		EntFriction:  engine.Engine.PMoveEntFriction(),
		EntGravity:   engine.Engine.PMoveEntGravity(),
		FrameTime:    engine.Engine.PMoveFrameTime(),
		Buttons:      engine.Engine.PMoveCmdButtons(),
		Onground:     engine.Engine.PMoveOnground(),
		Flags:        engine.Engine.PMoveFlags(),
		Waterlevel:   engine.Engine.PMoveWaterlevel(),
		InDuck:       engine.Engine.PMoveInDuck(),
		Impulse:      engine.Engine.PMoveImpulse(),
	})
	if err == nil {
		feed.Broadcast(binary)
	}

	gamelibs.Model.API().PMPlayerMove(server)
}

func (h *gameHandler) CGraphInitGraph(this unsafe.Pointer) {
	gamelibs.Model.API().CGraphInitGraph(this)
	engine.WorldGraph.SetPointer(this)
}
