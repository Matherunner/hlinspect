package hl

import "unsafe"

// SoundEnt represents CSoundEnt
type SoundEnt struct {
	address uintptr
}

// Sound represents CSound
type Sound struct {
	address uintptr
}

// Origin returns CSound::m_vecOrigin
func (sound *Sound) Origin() [3]float32 {
	return *(*[3]float32)(unsafe.Pointer(sound.address))
}

// Type returns CSound::m_iType
func (sound *Sound) Type() int32 {
	return *(*int32)(unsafe.Pointer(sound.address + 0xc))
}

// Volume returns CSound::m_iVolume
func (sound *Sound) Volume() int32 {
	return *(*int32)(unsafe.Pointer(sound.address + 0x10))
}

// ExpireTime returns CSound::m_flExpireTime
func (sound *Sound) ExpireTime() float32 {
	return *(*float32)(unsafe.Pointer(sound.address + 0x14))
}

// Next returns CSound::m_iNext
func (sound *Sound) Next() int32 {
	return *(*int32)(unsafe.Pointer(sound.address + 0x18))
}

// MakeSound creates an instance of Sound
func MakeSound(address uintptr) Sound {
	return Sound{address: address}
}

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
	soundIdx := CSoundEntActiveList()
	for soundIdx != -1 {
		address := CSoundEntSoundPointerForIndex(soundIdx)
		if address == 0 {
			break
		}
		sound := MakeSound(address)
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
