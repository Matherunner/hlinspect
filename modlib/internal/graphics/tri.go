package graphics

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs/hw"
)

// DrawTriangles draw OpenGL triangles
func DrawTriangles() {
	hw.TriGLColor4f(1, 0, 0, 0.2)
	hw.TriGLCullFace(hw.TriNone)
	hw.TriGLRenderMode(hw.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		className := engine.Engine.GlobalVariables.String(edict.EntVars().Classname())
		if className == "scripted_sequence" {
			origin := edict.EntVars().Origin()
			drawPyramid(origin, 10, 20)
		}
	}
}
