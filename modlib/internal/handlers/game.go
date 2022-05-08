package handlers

import (
	"context"
	"hlinspect/internal/cmd"
	"hlinspect/internal/cvar"
	"hlinspect/internal/feed"
	"hlinspect/internal/game"
	"hlinspect/internal/game/engine"
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
	game.Model.API().MemoryInit(buf, size)

	game.Model.API().RegisterCVar(&cvar.Wallhack)
	game.Model.API().RegisterCVar(&cvar.WallhackAlpha)
	game.Model.API().RegisterCVar(&cvar.FadeRemove)
	game.Model.API().RegisterCVar(&cvar.Nodes)
	game.Model.API().RegisterCVar(&cvar.NodeLinks)
	game.Model.API().RegisterCVar(&cvar.Cine)
	game.Model.API().RegisterCVar(&cvar.CinePossess)

	for name := range cmd.CommandHandlerByName {
		game.Model.API().CmdAddCommand(name)
	}
}

func (h *gameHandler) OnCommand() {
	name := game.Model.API().CmdArgv(0)
	if handler, ok := cmd.CommandHandlerByName[name]; ok {
		handler()
	}
}

func (h *gameHandler) VFadeAlpha() int {
	if cvar.FadeRemove.Float32() != 0 {
		return 0
	}
	return game.Model.API().VFadeAlpha()
}

func (h *gameHandler) RClear() {
	if cvar.Wallhack.Float32() != 0 {
		game.Model.GL().ClearColor(0, 0, 0, 1)
		game.Model.GL().Clear(game.GLColorBufferBit)
	}
	game.Model.API().RClear()
}

func (h *gameHandler) RDrawSequentialPoly(surf uintptr, free int) {
	if cvar.Wallhack.Float32() == 0 {
		game.Model.API().RDrawSequentialPoly(surf, free)
		return
	}

	game.Model.GL().Enable(game.GLBlend)
	game.Model.GL().DepthMask(false)
	game.Model.GL().BlendFunc(game.GLSrcAlpha, game.GLOneMinusSrcAlpha)
	game.Model.GL().Color4f(1, 1, 1, cvar.WallhackAlpha.Float32())

	game.Model.API().RDrawSequentialPoly(surf, free)

	game.Model.GL().DepthMask(true)
	game.Model.GL().Disable(game.GLBlend)
}

func (h *gameHandler) HUDRedraw(time float32, intermission int) {
	game.Model.API().HUDRedraw(time, intermission)

	graphics.DrawHUD(time, intermission)
}

func (h *gameHandler) HUDDrawTransparentTriangle() {
	game.Model.API().HUDDrawTransparentTriangle()

	game.Model.GL().Disable(game.GLTexture2D)
	graphics.DrawTriangles()
	game.Model.GL().Enable(game.GLTexture2D)
	game.Model.API().TriGLRenderMode(game.KRenderNormal)
}

func (h *gameHandler) HUDVidInit() int {
	ret := game.Model.API().HUDVidInit()

	screenInfo := game.Model.API().GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)

	return ret
}

func (h *gameHandler) HUDReset() {
	game.Model.API().HUDReset()

	screenInfo := game.Model.API().GetScreenInfo()
	graphics.SetScreenInfo(&screenInfo)
}

func (h *gameHandler) PMInit(ppm unsafe.Pointer) {
	game.Model.API().PMInit(ppm)

	game.Model.S().SetPPMove(ppm)
	logs.DLLLog.Debugf("Set PPMOVE with address = %x", ppm)
}

func (h *gameHandler) PMPlayerMove(server int) {
	binary, err := proto.Serialize(&proto.PMove{
		Stage:        proto.PMoveStagePre,
		Velocity:     game.Model.S().PMoveVelocity(),
		Position:     game.Model.S().PMovePosition(),
		Viewangles:   game.Model.S().PMoveViewangles(),
		Basevelocity: game.Model.S().PMoveBasevelocity(),
		FSU:          game.Model.S().PMoveCmdFSU(),
		Punchangles:  game.Model.S().PMovePunchangles(),
		EntFriction:  game.Model.S().PMoveEntFriction(),
		EntGravity:   game.Model.S().PMoveEntGravity(),
		FrameTime:    game.Model.S().PMoveFrameTime(),
		Buttons:      game.Model.S().PMoveCmdButtons(),
		Onground:     game.Model.S().PMoveOnground(),
		Flags:        game.Model.S().PMoveFlags(),
		Waterlevel:   game.Model.S().PMoveWaterlevel(),
		InDuck:       game.Model.S().PMoveInDuck(),
		Impulse:      game.Model.S().PMoveImpulse(),
	})
	if err == nil {
		feed.Broadcast(binary)
	}

	game.Model.API().PMPlayerMove(server)
}

func (h *gameHandler) CGraphInitGraph(this unsafe.Pointer) {
	game.Model.API().CGraphInitGraph(this)
	engine.WorldGraph.SetPointer(this)
}

func (h *gameHandler) CLCreateMove(frameTime float32, usercmd unsafe.Pointer, active int) {
	respChan := game.Model.Sync().InputControlResp(false)
	if respChan != nil {
		respChan <- true
	}

	reqChan := game.Model.Sync().InputControlReq(false)
	if reqChan != nil {
		cmd := <-reqChan
		logs.DLLLog.Debugf("CMD = %+v", cmd)
	}

	game.Model.API().CLCreateMove(frameTime, usercmd, active)
}
