package cvar

import (
	"hlinspect/internal/engine"
	"strconv"
	"unsafe"
)

import "C"

var Wallhack = makeCVar("hli_wallhack", "0")
var WallhackAlpha = makeCVar("hli_wallhack_alpha", "0.6")
var FadeRemove = makeCVar("hli_fade_remove", "0")
var Nodes = makeCVar("hli_nodes", "0")
var NodeLinks = makeCVar("hli_node_links", "1")
var Cine = makeCVar("hli_cine", "0")
var CinePossess = makeCVar("hli_cine_possess", "1")
var CollisionHullShow = makeCVar("hli_collision_hull_show", "0")

func makeCVar(name, value string) engine.CVar {
	floatVal, _ := strconv.ParseFloat(value, 32)
	cvar := engine.RawCVar{
		// Probably don't need to free these strings?
		Name:   uintptr(unsafe.Pointer(C.CString(name))),
		String: uintptr(unsafe.Pointer(C.CString(value))),
		Value:  float32(floatVal),
	}
	return engine.MakeCVar(unsafe.Pointer(&cvar))
}
