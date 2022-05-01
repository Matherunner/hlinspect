package gamelibs

// TODO: should this file really be here? or on a different layer?

import (
	"hlinspect/internal/engine"
	"unsafe"
)

// SoundItem is an internal representation of a sound
type SoundItem struct {
	Origin     [3]float32
	Type       int32
	Volume     int32
	ExpireTime float32
}

// GetSoundList returns a list of sounds in the game
func GetSoundList() []SoundItem {
	items := make([]SoundItem, 0, 10)
	soundIdx := Model.API().CSoundEntActiveList()
	for soundIdx != -1 {
		address := Model.API().CSoundEntSoundPointerForIndex(soundIdx)
		if address == nil {
			break
		}
		sound := engine.MakeSound(unsafe.Pointer(address))
		items = append(items, SoundItem{
			Origin:     sound.Origin(),
			Type:       sound.Type(),
			Volume:     sound.Volume(),
			ExpireTime: sound.ExpireTime(),
		})
		soundIdx = sound.Next()
	}
	return items
}
