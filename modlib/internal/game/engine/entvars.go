package engine

import "unsafe"

// EntVars represents entvars_t
type EntVars struct {
	ptr unsafe.Pointer
}

// MakeEntVars creates an instance of EntVars
func MakeEntVars(pointer unsafe.Pointer) EntVars {
	return EntVars{ptr: pointer}
}

// Ptr returns the pointer to this entvars_t object
func (entvars EntVars) Ptr() unsafe.Pointer {
	return entvars.ptr
}

// Origin returns entvars_t::origin
func (entvars EntVars) Origin() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0x8))
}

// Velocity returns entvars_t::velocity
func (entvars EntVars) Velocity() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 32))
}

// BaseVelocity returns envars_t::basevelocity
func (entvars EntVars) BaseVelocity() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 44))
}

// Angles returns entvars_t::angles
func (entvars EntVars) Angles() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0x50))
}

// PunchAngles returns entvars_t::punchangle
func (entvars EntVars) PunchAngles() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 104))
}

// EntityFriction returns entvars_t::friction
func (entvars EntVars) EntityFriction() float32 {
	return *(*float32)(unsafe.Add(entvars.ptr, 288))
}

// EntityGravity returns entvars_t::gravity
func (entvars EntVars) EntityGravity() float32 {
	return *(*float32)(unsafe.Add(entvars.ptr, 284))
}

// WaterLevel returns entvars_t::waterlevel
func (entvars EntVars) WaterLevel() int {
	return *(*int)(unsafe.Add(entvars.ptr, 448))
}

// GroundEntity returns entvars_t::groundentity
func (entvars EntVars) GroundEntity() Edict {
	p := *(*unsafe.Pointer)(unsafe.Add(entvars.ptr, 412))
	return MakeEdict(p)
}

// Flags returns entvars_t::flags
func (entvars EntVars) Flags() int {
	return *(*int)(unsafe.Add(entvars.ptr, 420))
}

// Classname returns entvars_t::classname
func (entvars EntVars) Classname() uint {
	return *(*uint)(unsafe.Add(entvars.ptr, 0x0))
}

// Targetname returns entvars_t::targetname
func (entvars EntVars) Targetname() uint {
	return *(*uint)(unsafe.Add(entvars.ptr, 0x1cc))
}

// Mins returns entvars_t::mins
func (entvars EntVars) Mins() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xdc))
}

// Maxs returns entvars_t::maxs
func (entvars EntVars) Maxs() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xe8))
}

// AbsMin returns entvars_t::absmin
func (entvars EntVars) AbsMin() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xc4))
}

// AbsMax returns entvars_t::absmax
func (entvars EntVars) AbsMax() [3]float32 {
	return *(*[3]float32)(unsafe.Add(entvars.ptr, 0xd0))
}
