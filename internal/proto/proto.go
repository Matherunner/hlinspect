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

// PMove represents the data held by `pmove`
type PMove struct {
	Velocity [3]float32
}

func (pmove *PMove) serialize() ([]byte, error) {
	v := &feed.PMove{Velocity: pmove.Velocity[:]}
	return proto.Marshal(v)
}

// Serialize converts the given `s` structure into bytes
func Serialize(s serializer) ([]byte, error) {
	return s.serialize()
}
