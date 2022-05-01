package cmd

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs"
	"math"
	"unsafe"
)

var TrackedNPC = map[unsafe.Pointer]bool{}
var ShowRadiusCine = map[unsafe.Pointer]bool{}

// To add a new command, simply specify the name and the function here. We don't need
// to edit anywhere else.
var CommandHandlerByName = map[string]func(){
	"hli_npc_track_add": func() {
		// TODO: uncomment this

		position := engine.Engine.PMovePosition()
		position[2] += 28

		viewangles := engine.Engine.PMoveViewangles()

		forward, _, _ := gamelibs.Model.API().AngleVectors(viewangles)
		endPos := [3]float32{}
		for i := 0; i < 3; i++ {
			// TODO: debug temporary
			position[i] += 20 * forward[i]
			endPos[i] = position[i] + 8192*forward[i]
		}

		result := gamelibs.Model.API().TraceLine(position, endPos, gamelibs.TraceDontIgnoreMonsters, nil)
		edict := engine.MakeEdict(result.Hit)
		if edict.Pointer() == nil || edict.Free() || edict.PrivateData() == nil {
			return
		}

		// FIXME: will crash in HUD if not a monster

		TrackedNPC[edict.PrivateData()] = true
	},
	"hli_npc_track_del": func() {
		// TODO: check for * and param, delete all for now
		TrackedNPC = map[unsafe.Pointer]bool{}
	},
	"hli_npc_track_list": func() {
		// TODO: maybe print out the class names
	},
	"hli_cine_radius_all": func() {

	},
	"hli_cine_radius_nearest_add": func() {
		position := engine.Engine.PMovePosition()
		minDistance := math.MaxFloat64
		var minEnt unsafe.Pointer
		numEdicts := engine.Engine.SV.NumEdicts()
		for i := 0; i < numEdicts; i++ {
			edict := engine.Engine.SV.Edict(i)
			if edict.Free() {
				continue
			}

			entVars := edict.EntVars()
			className := engine.Engine.GlobalVariables.String(entVars.Classname())
			if className != "scripted_sequence" {
				continue
			}

			entOrigin := entVars.Origin()
			distance := math.Hypot(float64(entOrigin[0]-position[0]), float64(entOrigin[1]-position[1]))
			if distance < minDistance {
				minDistance = distance
				minEnt = edict.PrivateData()
			}
		}
		ShowRadiusCine[minEnt] = true
	},
	"hli_cine_radius_nearest_del": func() {

	},
	"hli_cine_radius_clear": func() {
		ShowRadiusCine = map[unsafe.Pointer]bool{}
	},
}
