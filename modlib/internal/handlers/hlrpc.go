package handlers

import (
	"context"
	"hlinspect/internal/engine"
	"hlinspect/internal/game"
	"hlinspect/internal/hlrpc/schema"
	"hlinspect/internal/logs"
)

type HLRPCHandler struct {
}

func NewHLRPCHandler() *HLRPCHandler {
	return &HLRPCHandler{}
}

func (h *HLRPCHandler) GetFullPlayerState(ctx context.Context, resp *schema.FullPlayerState) error {
	pos := engine.Engine.PMovePosition()
	resp.SetPositionX(pos[0])
	resp.SetPositionY(pos[1])
	resp.SetPositionZ(pos[2])

	vel := engine.Engine.PMoveVelocity()
	resp.SetVelocityX(vel[0])
	resp.SetVelocityY(vel[1])
	resp.SetVelocityZ(vel[2])

	basevel := engine.Engine.PMoveBasevelocity()
	resp.SetBaseVelocityX(basevel[0])
	resp.SetBaseVelocityY(basevel[1])
	resp.SetBaseVelocityZ(basevel[2])

	angles := engine.Engine.PMoveViewangles()
	resp.SetYaw(angles[1])
	resp.SetPitch(angles[0])
	resp.SetRoll(angles[2])

	punchangles := engine.Engine.PMovePunchangles()
	resp.SetPunchYaw(punchangles[1])
	resp.SetPunchPitch(punchangles[0])
	resp.SetPunchRoll(punchangles[2])

	resp.SetEntityFriction(engine.Engine.PMoveEntFriction())
	resp.SetEntityGravity(engine.Engine.PMoveEntGravity())

	resp.SetOnGround(engine.Engine.PMoveOnground())
	// TODO: duck state
	resp.SetWaterLevel(uint8(engine.Engine.PMoveWaterlevel()))

	return nil
}

func (h *HLRPCHandler) StartInputControl(ctx context.Context) error {
	game.Model.Sync().StartInputControl()

	return nil
}

func (h *HLRPCHandler) StopInputControl(ctx context.Context) error {
	game.Model.Sync().StopInputControl()

	return nil
}

func (h *HLRPCHandler) InputStep(ctx context.Context, req *schema.CommandInput, resp *schema.FullPlayerState) error {
	reqChan := game.Model.Sync().InputControlReq(true)
	if reqChan == nil {
		// TODO: maybe return an error?
		return nil
	}

	logs.DLLLog.Debugf("Serve waiting to write req")
	reqChan <- req

	respChan := game.Model.Sync().InputControlResp(true)
	if respChan == nil {
		// TODO: return an error?
		return nil
	}

	logs.DLLLog.Debugf("Serve waiting for response")
	ret := <-respChan
	logs.DLLLog.Debugf("Serve GOT RESP")
	if ret {
		// TODO: return current state
		pos := engine.Engine.PMovePosition()
		resp.SetPositionX(pos[0])
		resp.SetPositionY(pos[1])
		resp.SetPositionZ(pos[2])

		vel := engine.Engine.PMoveVelocity()
		resp.SetVelocityX(vel[0])
		resp.SetVelocityY(vel[1])
		resp.SetVelocityZ(vel[2])

		basevel := engine.Engine.PMoveBasevelocity()
		resp.SetBaseVelocityX(basevel[0])
		resp.SetBaseVelocityY(basevel[1])
		resp.SetBaseVelocityZ(basevel[2])

		angles := engine.Engine.PMoveViewangles()
		resp.SetYaw(angles[1])
		resp.SetPitch(angles[0])
		resp.SetRoll(angles[2])

		punchangles := engine.Engine.PMovePunchangles()
		resp.SetPunchYaw(punchangles[1])
		resp.SetPunchPitch(punchangles[0])
		resp.SetPunchRoll(punchangles[2])

		resp.SetEntityFriction(engine.Engine.PMoveEntFriction())
		resp.SetEntityGravity(engine.Engine.PMoveEntGravity())

		resp.SetOnGround(engine.Engine.PMoveOnground())
		// TODO: duck state
		resp.SetWaterLevel(uint8(engine.Engine.PMoveWaterlevel()))
	}

	return nil
}
