// Package proto provides a translation layer to decouple serialisation and game logic.
package proto

import (
	"hlinspect/internal/proto/feed/feed"

	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=. feed.proto

type serializer interface {
	serialize() ([]byte, error)
}

const (
	PMoveStagePre  = int32(1)
	PMoveStagePost = int32(2)
)

// PMove represents the data held by `pmove`
type PMove struct {
	Stage        int32
	Velocity     [3]float32
	Position     [3]float32
	Viewangles   [3]float32
	Basevelocity [3]float32
	FSU          [3]float32
	Punchangles  [3]float32
	EntFriction  float32
	EntGravity   float32
	FrameTime    uint32
	Buttons      uint32
	Onground     bool
	Flags        uint32
	Waterlevel   uint32
	InDuck       bool
	Ladder       bool
	Impulse      uint32
}

func (pmove *PMove) serialize() ([]byte, error) {
	return proto.Marshal(&feed.PMove{
		Velocity:     pmove.Velocity[:],
		Position:     pmove.Position[:],
		Viewangles:   pmove.Viewangles[:],
		Basevelocity: pmove.Basevelocity[:],
		Fsu:          pmove.FSU[:],
		Punchangles:  pmove.Punchangles[:],
		EntFriction:  pmove.EntFriction,
		EntGravity:   pmove.EntGravity,
		FrameTime:    pmove.FrameTime,
		Buttons:      pmove.Buttons,
		Onground:     pmove.Onground,
		Flags:        pmove.Flags,
		Waterlevel:   pmove.Waterlevel,
		InDuck:       pmove.InDuck,
		Ladder:       pmove.Ladder,
		Impulse:      pmove.Impulse,
	})
}

// Serialize converts the given `s` structure into bytes
func Serialize(s serializer) ([]byte, error) {
	return s.serialize()
}
