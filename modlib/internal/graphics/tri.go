package graphics

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs/hw"
)

// DrawTriangles draw OpenGL triangles
func DrawTriangles() {
	drawScriptedSequences()
	drawScriptedSequencesPossessions()
}

func drawScriptedSequences() {
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

			cine := engine.MakeCine(edict.PrivateData())
			radius := cine.Radius()
			drawSphere(origin, radius, 10, 10)
		}
	}
}

func drawScriptedSequencesPossessions() {
	GLLineWidth(4)
	hw.TriGLColor4f(1, 0, 0, 1)
	hw.TriGLCullFace(hw.TriNone)
	hw.TriGLRenderMode(hw.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		className := engine.Engine.GlobalVariables.String(edict.EntVars().Classname())
		if className == "monster_human_torch_ally" || className == "monster_scientist" {
			cineAddr := engine.MakeMonster(edict.PrivateData()).CineAddr()
			if cineAddr == 0 {
				continue
			}

			entVars := engine.MakeEntity(cineAddr).EntVars()
			cineOrigin := entVars.Origin()
			monsterOrigin := edict.EntVars().Origin()
			drawLines([][3]float32{cineOrigin, monsterOrigin})
		}
	}
}
