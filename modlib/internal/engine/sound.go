package engine

import "unsafe"

const (
	SoundNone    = 0
	SoundCombat  = 1 << 0
	SoundWorld   = 1 << 1
	SoundPlayer  = 1 << 2
	SoundCarcass = 1 << 3
	SoundMeat    = 1 << 4
	SoundDanger  = 1 << 5
	SoundGarbage = 1 << 6
)

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
