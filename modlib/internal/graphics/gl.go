package graphics

/*
#cgo 386 LDFLAGS: -lopengl32

#include <GL/gl.h>
*/
import "C"

const (
	// GLTexture2D GL_TEXTURE_2D
	GLTexture2D = C.GL_TEXTURE_2D
)

// GLEnable glEnable
func GLEnable(cap uint) {
	C.glEnable(C.uint(cap))
}

// GLDisable glDisable
func GLDisable(cap uint) {
	C.glDisable(C.uint(cap))
}

// GLLineWidth glLineWidth
func GLLineWidth(width float32) {
	C.glLineWidth(C.float(width))
}
