package handlers

import (
	"context"
	"hlinspect/internal/engine"
	"hlinspect/internal/hlrpc/schema"
)

type hlrpcHandler struct {
}

func NewHLRPCHandler() *hlrpcHandler {
	return &hlrpcHandler{}
}

func (h *hlrpcHandler) GetFullPlayerState(ctx context.Context, resp *schema.FullPlayerState) error {
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
