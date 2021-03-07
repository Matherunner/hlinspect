package graphics

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs/hw"
	"hlinspect/internal/gl"
)

// DrawTriangles draw OpenGL triangles
func DrawTriangles() {
	drawNodeGraph()
	drawNodeLinks()
	drawScriptedSequences()
	drawScriptedSequencesPossessions()
	drawMonsterRoutes()
	drawBoundingBoxes()
}

func drawNodeGraph() {
	hw.TriGLColor4f(1, 0, 1, 0.2)
	hw.TriGLCullFace(hw.TriNone)
	hw.TriGLRenderMode(hw.KRenderTransAdd)

	numNodes := engine.WorldGraph.NumNodes()
	for i := 0; i < numNodes; i++ {
		node := engine.WorldGraph.Node(i)
		origin := node.Origin()
		drawPyramid(origin, 10, 20)
	}
}

func drawNodeLinks() {
	hw.TriGLColor4f(1, 0, 1, 0.2)
	hw.TriGLCullFace(hw.TriNone)
	hw.TriGLRenderMode(hw.KRenderTransAdd)

	numLinks := engine.WorldGraph.NumLinks()
	for i := 0; i < numLinks; i++ {
		link := engine.WorldGraph.Link(i)
		src := link.Source().Origin()
		dest := link.Destination().Origin()
		drawLines([][3]float32{src, dest})

		entvars := link.LinkEnt()
		if entvars.Pointer() != nil {
			origin := entvars.Origin()
			mins := entvars.Mins()
			maxs := entvars.Maxs()
			for i := 0; i < 3; i++ {
				mins[i] += origin[i]
				maxs[i] += origin[i]
			}
			// TODO: maybe tone down the brightness for this
			drawAACuboid(mins, maxs)
		}
	}
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

			// TODO: re-enable this?
			// cine := engine.MakeCine(edict.PrivateData())
			// radius := cine.Radius()
			// drawSphere(origin, radius, 10, 10)
		}
	}
}

func drawScriptedSequencesPossessions() {
	gl.LineWidth(4)
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
			cine := engine.MakeMonster(edict.PrivateData()).Cine()
			if cine.Pointer() == nil {
				continue
			}

			entity := engine.MakeEntity(cine.Pointer())
			cineOrigin := entity.EntVars().Origin()
			monsterOrigin := edict.EntVars().Origin()
			drawLines([][3]float32{cineOrigin, monsterOrigin})
		}
	}
}

func drawMonsterRoutes() {

}

func drawBoundingBoxes() {
	gl.LineWidth(1)
	hw.TriGLColor4f(0, 1, 0, 1)
	hw.TriGLCullFace(hw.TriNone)
	hw.TriGLRenderMode(hw.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		// Mins and Maxs are more accurate than AbsMin and AbsMax, see alien grunt's bbox
		entVars := edict.EntVars()
		className := engine.Engine.GlobalVariables.String(entVars.Classname())
		if className == "monster_alien_grunt" {
			origin := entVars.Origin()
			mins := entVars.Mins()
			maxs := entVars.Maxs()
			for i := 0; i < 3; i++ {
				mins[i] += origin[i]
				maxs[i] += origin[i]
			}
			drawAACuboidWireframe(mins, maxs)
		}
	}
}
