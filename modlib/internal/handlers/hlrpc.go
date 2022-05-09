package handlers

import (
	"context"
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
	ev := game.Model.S().SVPlayer().EntVars()

	pos := ev.Origin()
	resp.SetPositionX(pos[0])
	resp.SetPositionY(pos[1])
	resp.SetPositionZ(pos[2])

	vel := ev.Velocity()
	resp.SetVelocityX(vel[0])
	resp.SetVelocityY(vel[1])
	resp.SetVelocityZ(vel[2])

	basevel := ev.BaseVelocity()
	resp.SetBaseVelocityX(basevel[0])
	resp.SetBaseVelocityY(basevel[1])
	resp.SetBaseVelocityZ(basevel[2])

	angles := ev.Angles()
	resp.SetYaw(angles[1])
	resp.SetPitch(angles[0])
	resp.SetRoll(angles[2])

	punchangles := ev.PunchAngles()
	resp.SetPunchYaw(punchangles[1])
	resp.SetPunchPitch(punchangles[0])
	resp.SetPunchRoll(punchangles[2])

	resp.SetEntityFriction(ev.EntityFriction())
	resp.SetEntityGravity(ev.EntityGravity())

	resp.SetOnGround((ev.Flags() & game.FLOnGround) != 0)
	// TODO: duck state
	resp.SetWaterLevel(uint8(ev.WaterLevel()))

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
		pos := game.Model.S().PMovePosition()
		resp.SetPositionX(pos[0])
		resp.SetPositionY(pos[1])
		resp.SetPositionZ(pos[2])

		vel := game.Model.S().PMoveVelocity()
		resp.SetVelocityX(vel[0])
		resp.SetVelocityY(vel[1])
		resp.SetVelocityZ(vel[2])

		basevel := game.Model.S().PMoveBasevelocity()
		resp.SetBaseVelocityX(basevel[0])
		resp.SetBaseVelocityY(basevel[1])
		resp.SetBaseVelocityZ(basevel[2])

		angles := game.Model.S().PMoveViewangles()
		resp.SetYaw(angles[1])
		resp.SetPitch(angles[0])
		resp.SetRoll(angles[2])

		punchangles := game.Model.S().PMovePunchangles()
		resp.SetPunchYaw(punchangles[1])
		resp.SetPunchPitch(punchangles[0])
		resp.SetPunchRoll(punchangles[2])

		resp.SetEntityFriction(game.Model.S().PMoveEntFriction())
		resp.SetEntityGravity(game.Model.S().PMoveEntGravity())

		resp.SetOnGround(game.Model.S().PMoveOnground())
		// TODO: duck state
		resp.SetWaterLevel(uint8(game.Model.S().PMoveWaterlevel()))
	}

	return nil
}
