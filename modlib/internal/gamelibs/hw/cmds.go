package hw

import (
	"hlinspect/internal/engine"
	"unsafe"
)

var TrackedNPC = map[unsafe.Pointer]bool{}

var commandHandlerByName = map[string]func(){
	"hli_npc_track_add": func() {
		position := engine.Engine.PMovePosition()
		position[2] += 28

		viewangles := engine.Engine.PMoveViewangles()

		forward, _, _ := AngleVectors(viewangles)
		endPos := [3]float32{}
		for i := 0; i < 3; i++ {
			// TODO: debug temporary
			position[i] += 20 * forward[i]
			endPos[i] = position[i] + 8192*forward[i]
		}

		result := TraceLine(position, endPos, TraceDontIgnoreMonsters, nil)
		edict := engine.MakeEdict(result.Hit)
		if edict.Pointer() == nil || edict.Free() || edict.PrivateData() == nil {
			return
		}

		// FIXME: will crash in HUD if not a monster

		TrackedNPC[edict.PrivateData()] = true
	},

	"hli_npc_track_del": func() {

	},

	"hli_npc_track_list": func() {

	},
}
