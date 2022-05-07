package graphics

import (
	"hlinspect/internal/game"
	"math"
)

func drawAACuboid(corner1, corner2 [3]float32) {
	game.Model.API().TriGLBegin(game.TriQuads)

	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})

	game.Model.API().TriGLEnd()
}

func drawAACuboidWireframe(corner1, corner2 [3]float32) {
	game.Model.API().TriGLBegin(game.TriLines)

	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner1[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner1[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})

	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner1[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner1[0], corner2[1], corner2[2]})
	game.Model.API().TriGLVertex3fv([3]float32{corner2[0], corner2[1], corner2[2]})

	game.Model.API().TriGLEnd()
}

func drawPyramid(origin [3]float32, width, height float32) {
	halfWidth := width * 0.5
	offsets := [5][2]float32{
		{halfWidth, halfWidth},
		{halfWidth, -halfWidth},
		{-halfWidth, -halfWidth},
		{-halfWidth, halfWidth},
		{halfWidth, halfWidth},
	}

	game.Model.API().TriGLBegin(game.TriQuads)
	for _, offset := range offsets {
		game.Model.API().TriGLVertex3fv([3]float32{origin[0] + offset[0], origin[1] + offset[1], origin[2]})
	}
	game.Model.API().TriGLEnd()

	game.Model.API().TriGLBegin(game.TriTriangleFan)
	game.Model.API().TriGLVertex3fv([3]float32{origin[0], origin[1], origin[2] + height})
	for _, offset := range offsets {
		game.Model.API().TriGLVertex3fv([3]float32{origin[0] + offset[0], origin[1] + offset[1], origin[2]})
	}
	game.Model.API().TriGLEnd()
}

func drawInvertedPyramid(origin [3]float32, width, height float32) {
	halfWidth := width * 0.5
	offsets := [5][2]float32{
		{halfWidth, halfWidth},
		{halfWidth, -halfWidth},
		{-halfWidth, -halfWidth},
		{-halfWidth, halfWidth},
		{halfWidth, halfWidth},
	}

	game.Model.API().TriGLBegin(game.TriQuads)
	for _, offset := range offsets {
		game.Model.API().TriGLVertex3fv([3]float32{origin[0] + offset[0], origin[1] + offset[1], origin[2] + height})
	}
	game.Model.API().TriGLEnd()

	game.Model.API().TriGLBegin(game.TriTriangleFan)
	game.Model.API().TriGLVertex3fv([3]float32{origin[0], origin[1], origin[2]})
	for _, offset := range offsets {
		game.Model.API().TriGLVertex3fv([3]float32{origin[0] + offset[0], origin[1] + offset[1], origin[2] + height})
	}
	game.Model.API().TriGLEnd()
}

func drawLines(positions [][3]float32) {
	game.Model.API().TriGLBegin(game.TriLines)
	for i := 0; i < len(positions)-1; i++ {
		game.Model.API().TriGLVertex3fv(positions[i])
		game.Model.API().TriGLVertex3fv(positions[i+1])
	}
	game.Model.API().TriGLEnd()
}

func drawSphere(origin [3]float32, r float32, nlat int, nlong int) {
	for i := 0; i <= nlat; i++ {
		lat0 := (-0.5 + float64(i-1)/float64(nlat)) * math.Pi
		z0, zr0 := math.Sincos(lat0)

		lat1 := (-0.5 + float64(i)/float64(nlat)) * math.Pi
		z1, zr1 := math.Sincos(lat1)

		game.Model.API().TriGLBegin(game.TriQuadStrip)
		for j := 0; j <= nlong; j++ {
			lng := float64(j-1) / float64(nlong) * 2.0 * math.Pi
			y, x := math.Sincos(lng)
			game.Model.API().TriGLVertex3fv([3]float32{
				r*float32(x*zr0) + origin[0],
				r*float32(y*zr0) + origin[1],
				r*float32(z0) + origin[2],
			})
			game.Model.API().TriGLVertex3fv([3]float32{
				r*float32(x*zr1) + origin[0],
				r*float32(y*zr1) + origin[1],
				r*float32(z1) + origin[2],
			})
		}
		game.Model.API().TriGLEnd()
	}
}

func worldToHUDScreen(point [3]float32, width, height int) (screen [2]int, clipped bool) {
	fscreen, clipped := game.Model.API().ScreenTransform(point)
	screen[0] = int((1 + fscreen[0]) * 0.5 * float32(width))
	screen[1] = int((1 - fscreen[1]) * 0.5 * float32(height))
	return
}
