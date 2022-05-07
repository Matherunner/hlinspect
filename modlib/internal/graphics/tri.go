package graphics

import (
	"hlinspect/internal/cmd"
	"hlinspect/internal/cvar"
	"hlinspect/internal/engine"
	"hlinspect/internal/game"
	"strings"
)

// DrawTriangles draw OpenGL triangles
func DrawTriangles() {
	drawNodeGraph()
	drawNodeLinks()
	drawScriptedSequences()
	drawScriptedSequencesPossessions()
	drawMonsterRoutes()
	drawBoundingBoxes()
	drawSoundLinks()
	drawInfoBigMomma()
}

func drawNodeGraph() {
	if cvar.Nodes.Float32() == 0 {
		return
	}

	game.Model.API().TriGLColor4f(1, 0, 1, 0.2)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

	numNodes := engine.WorldGraph.NumNodes()
	for i := 0; i < numNodes; i++ {
		node := engine.WorldGraph.Node(i)
		origin := node.Origin()
		drawPyramid(origin, 10, 20)
	}
}

func drawNodeLinks() {
	if cvar.Nodes.Float32() == 0 || cvar.NodeLinks.Float32() == 0 {
		return
	}

	game.Model.API().TriGLColor4f(1, 0, 1, 0.2)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

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
	game.Model.API().TriGLColor4f(1, 0, 0, 0.2)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		className := engine.Engine.GlobalVariables.String(edict.EntVars().Classname())
		if className == "scripted_sequence" {
			cine := engine.MakeCine(edict.PrivateData())
			interruptible := cine.Interruptible()
			origin := edict.EntVars().Origin()
			if interruptible {
				drawInvertedPyramid(origin, 10, 20)
			} else {
				drawPyramid(origin, 10, 20)
			}

			if cmd.ShowRadiusCine[edict.PrivateData()] {
				radius := cine.Radius()
				drawSphere(origin, radius, 50, 50)
			}
		}
	}
}

func drawScriptedSequencesPossessions() {
	game.Model.GL().LineWidth(4)
	game.Model.API().TriGLColor4f(1, 0, 0, 1)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		className := engine.Engine.GlobalVariables.String(edict.EntVars().Classname())
		if className == "monster_human_torch_ally" || className == "monster_scientist" {
			monster := engine.MakeMonster(edict.PrivateData())
			cine := monster.Cine()
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
	game.Model.GL().LineWidth(4)
	game.Model.API().TriGLColor4f(0, 1, 0, 1)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		className := engine.Engine.GlobalVariables.String(edict.EntVars().Classname())
		if strings.HasPrefix(className, "monster_") {
			if edict.PrivateData() == nil {
				continue
			}

			monster := engine.MakeMonster(edict.PrivateData())
			routes := monster.Routes()
			routeIndex := monster.RouteIndex()
			monsterOrigin := edict.EntVars().Origin()
			lines := [][3]float32{monsterOrigin}
			for i := routeIndex; i < len(routes); i++ {
				route := routes[i]
				if route.Type() == 0 {
					break
				}
				lines = append(lines, route.Location())
				if route.Type()&engine.RouteMFIsGoal != 0 {
					break
				}
			}
			drawLines(lines)
		}
	}
}

func drawBoundingBoxes() {
	game.Model.GL().LineWidth(1)
	game.Model.API().TriGLColor4f(0, 1, 0, 1)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

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

func drawSoundLinks() {
	game.Model.GL().LineWidth(3)
	game.Model.API().TriGLColor4f(0.5, 0.8, 1, 1)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		entVars := edict.EntVars()
		className := engine.Engine.GlobalVariables.String(entVars.Classname())
		if strings.HasPrefix(className, "monster_") {
			origin := entVars.Origin()
			monster := engine.MakeMonster(edict.PrivateData())
			if monster.MonsterState() == engine.MonsterStateDead || monster.MonsterState() == engine.MonsterStateNone || monster.MonsterState() == engine.MonsterStateProne {
				continue
			}

			e := game.Model.API().PFCheckClientI(edict.Pointer())
			if e == 0 || engine.Engine.SV.EntOffset(e) == 0 {
				// Not in PVS
				if monster.MonsterState() != engine.MonsterStateCombat {
					// If this condition is true, then Listen is not called in the game
					// If we don't do this check, there will be a ton of false positives because m_iAudibleList is initialised to 0
					// which happens to be the player's sound.
					continue
				}
			}

			audibleList := monster.AudibleList()
			if audibleList == -1 {
				continue
			}

			sound := engine.MakeSound(game.Model.API().CSoundEntSoundPointerForIndex(int32(audibleList)))
			soundOrigin := sound.Origin()
			drawLines([][3]float32{origin, soundOrigin})
		}
	}
}

func drawInfoBigMomma() {
	game.Model.API().TriGLColor4f(0.8, 0.4, 0.7, 1)
	game.Model.API().TriGLCullFace(game.TriNone)
	game.Model.API().TriGLRenderMode(game.KRenderTransAdd)

	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		entVars := edict.EntVars()
		className := engine.Engine.GlobalVariables.String(entVars.Classname())
		if strings.HasPrefix(className, "info_bigmomma") {
			origin := entVars.Origin()
			drawPyramid(origin, 50, 100)
		}
	}
}
