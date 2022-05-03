package events

import (
	"context"
	"hlinspect/internal/cmd"
	"hlinspect/internal/cvar"
	"hlinspect/internal/engine"
	"hlinspect/internal/feed"
	"hlinspect/internal/gamelibs"
	"hlinspect/internal/gamelibs/gl"
	"hlinspect/internal/graphics"
	"hlinspect/internal/logs"
	"hlinspect/internal/proto"
	"unsafe"
)

type handler struct {
}

func NewHandler() gamelibs.Handler {
	return &handler{}
}

func (h *handler) MemoryInit(ctx context.Context, buf uintptr, size int) {
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

func (h *handler) PMInit(ppm unsafe.Pointer) {
	gamelibs.Model.API().PMInit(ppm)

	engine.Engine.SetPPMove(ppm)
	logs.DLLLog.Debugf("Set PPMOVE with address = %x", ppm)
}

func (h *handler) PMPlayerMove(server int) {
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

func (h *handler) CGraphInitGraph(this unsafe.Pointer) {
	gamelibs.Model.API().CGraphInitGraph(this)
	engine.WorldGraph.SetPointer(this)
}
