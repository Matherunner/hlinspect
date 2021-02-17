package graphics

import (
	"hlinspect/internal/gamelibs/hw"
)

// DrawHUD draw HUD
func DrawHUD(time float32, intermission int32) {
	hw.VGUI2DrawSetTextColorAlpha(255, 180, 30, 255)
	hw.DrawString(10, 10, "Hello world!")
}
