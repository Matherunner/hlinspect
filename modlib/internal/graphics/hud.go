package graphics

import (
	"fmt"
	"hlinspect/internal/cmd"
	"hlinspect/internal/game"
	"hlinspect/internal/game/engine"
	"strings"
)

var screenInfo *game.ScreenInfo

// SetScreenInfo sets the current screen info
func SetScreenInfo(si *game.ScreenInfo) {
	screenInfo = si
}

// DrawHUD draws HUD
func DrawHUD(time float32, intermission int) {
	game.Model.API().VGUI2DrawSetTextColorAlpha(255, 180, 30, 255)

	drawEntitiesOverlay()
	drawSounds()
}

func drawSounds() {
	sounds := game.GetSoundList()
	for _, sound := range sounds {
		screen, clipped := worldToHUDScreen(sound.Origin, int(screenInfo.Width), int(screenInfo.Height))
		if !clipped {
			game.Model.API().DrawString(screen[0], screen[1], fmt.Sprintf("%v", sound.Type))
		}
	}
}

func drawEntitiesOverlay() {
	numEdicts := game.Model.S().SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := game.Model.S().SV.Edict(i)
		if edict.Free() {
			continue
		}

		tracked := cmd.TrackedNPC[edict.PrivateData()]
		if !tracked {
			continue
		}

		entVars := edict.EntVars()
		if !strings.HasPrefix(game.Model.S().GlobalVariables.String(entVars.Classname()), "monster_") {
			continue
		}

		monster := engine.MakeMonster(edict.PrivateData())

		// TODO: commented out sound
		// {
		// 	origin := entVars.Origin()
		// 	screen, clipped := worldToHUDScreen(origin, int(screenInfo.Width), int(screenInfo.Height))
		// 	if !clipped {
		// 		nextAudible := monster.NextAudible()
		// 		hw.DrawString(screen[0], screen[1]-int(screenInfo.CharHeight), fmt.Sprintf("Audible: %v", nextAudible))
		// 	}
		// }

		schedule := monster.Schedule()
		if schedule != nil {
			origin := entVars.Origin()
			screen, clipped := worldToHUDScreen(origin, int(screenInfo.Width), int(screenInfo.Height))
			if !clipped {
				game.Model.API().VGUI2DrawSetTextColorAlpha(0, 255, 0, 255)
				game.Model.API().DrawString(screen[0], screen[1], fmt.Sprintf("State: %v", engine.MonsterStateToString(monster.MonsterState())))
				game.Model.API().VGUI2DrawSetTextColorAlpha(255, 180, 30, 255)
				game.Model.API().DrawString(screen[0], screen[1]+int(screenInfo.CharHeight), fmt.Sprintf("Schedule: %v", schedule.Name()))
				task := schedule.Task(monster.ScheduleIndex())
				game.Model.API().DrawString(screen[0], screen[1]+2*int(screenInfo.CharHeight), fmt.Sprintf("Task: %v (%v)", task.Name(), task.Data))
				angles := entVars.Angles()
				game.Model.API().DrawString(screen[0], screen[1]+3*int(screenInfo.CharHeight), fmt.Sprintf("%v %v", angles[0], angles[1]))
				game.Model.API().DrawString(screen[0], screen[1]+4*int(screenInfo.CharHeight), fmt.Sprintf("%v %v %v", origin[0], origin[1], origin[2]))

				cine := engine.MakeMonster(edict.PrivateData()).Cine()
				if cine.Pointer() != nil {
					if cine.Interruptible() {
						game.Model.API().DrawString(screen[0], screen[1]+5*int(screenInfo.CharHeight), "Interruptible")
					} else {
						game.Model.API().DrawString(screen[0], screen[1]+5*int(screenInfo.CharHeight), "Uninterruptible")
					}
				}

				game.Model.API().VGUI2DrawSetTextColorAlpha(255, 255, 0, 255)
				e := game.Model.API().PFCheckClientI(edict.Pointer())
				if e != 0 && game.Model.S().SV.EntOffset(e) != 0 {
					game.Model.API().DrawString(screen[0], screen[1]+6*int(screenInfo.CharHeight), "In PVS")
				} else {
					game.Model.API().DrawString(screen[0], screen[1]+6*int(screenInfo.CharHeight), "Not in PVS")
				}
				game.Model.API().VGUI2DrawSetTextColorAlpha(255, 180, 30, 255)
			}
		}
	}
}
