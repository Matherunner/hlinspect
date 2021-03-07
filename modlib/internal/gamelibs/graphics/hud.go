package graphics

import (
	"fmt"
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs/hl"
	"hlinspect/internal/gamelibs/hw"
)

var screenInfo *hw.ScreenInfo

// SetScreenInfo sets the current screen info
func SetScreenInfo(si *hw.ScreenInfo) {
	screenInfo = si
}

// DrawHUD draws HUD
func DrawHUD(time float32, intermission int32) {
	hw.VGUI2DrawSetTextColorAlpha(255, 180, 30, 255)

	drawEntitiesOverlay()
	drawSounds()
}

func drawSounds() {
	sounds := hl.GetSoundList()
	for _, sound := range sounds {
		screen, clipped := worldToHUDScreen(sound.Origin, int(screenInfo.Width), int(screenInfo.Height))
		if !clipped {
			hw.DrawString(screen[0], screen[1], fmt.Sprintf("%v", sound.Type))
		}
	}
}

func drawEntitiesOverlay() {
	numEdicts := engine.Engine.SV.NumEdicts()
	for i := 0; i < numEdicts; i++ {
		edict := engine.Engine.SV.Edict(i)
		if edict.Free() {
			continue
		}

		className := engine.Engine.GlobalVariables.String(edict.EntVars().Classname())
		if className == "monster_scientist" || className == "monster_barney" || className == "monster_bigmomma" || className == "monster_human_torch_ally" {
			monster := engine.MakeMonster(edict.PrivateData())
			schedule := monster.Schedule()
			if schedule != nil {
				origin := edict.EntVars().Origin()
				screen, clipped := worldToHUDScreen(origin, int(screenInfo.Width), int(screenInfo.Height))
				if !clipped {
					hw.DrawString(screen[0], screen[1], schedule.Name())
					task := schedule.Task(monster.ScheduleIndex())
					hw.DrawString(screen[0], screen[1]+int(screenInfo.CharHeight), fmt.Sprintf("%v (%v)", task.Name(), task.Data))
					angles := edict.EntVars().Angles()
					hw.DrawString(screen[0], screen[1]+2*int(screenInfo.CharHeight), fmt.Sprintf("%v %v", angles[0], angles[1]))
					hw.DrawString(screen[0], screen[1]+3*int(screenInfo.CharHeight), fmt.Sprintf("%v %v %v", origin[0], origin[1], origin[2]))

					cine := engine.MakeMonster(edict.PrivateData()).Cine()
					if cine.Pointer() != nil {
						if cine.Interruptible() {
							hw.DrawString(screen[0], screen[1]+4*int(screenInfo.CharHeight), "I")
						} else {
							hw.DrawString(screen[0], screen[1]+4*int(screenInfo.CharHeight), "UI")
						}
					}

					e := hw.PFCheckClientI(edict.Pointer())
					if e != 0 && engine.Engine.SV.EntOffset(e) != 0 {
						hw.DrawString(screen[0], screen[1]+5*int(screenInfo.CharHeight), "V")
					} else {
						hw.DrawString(screen[0], screen[1]+5*int(screenInfo.CharHeight), "IV")
					}
				}
			}
		}
	}
}
