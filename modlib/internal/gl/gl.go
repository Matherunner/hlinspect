package gl

/*
#cgo 386 LDFLAGS: -lopengl32

#include <GL/gl.h>
*/
import "C"

const (
	// Texture2D GL_TEXTURE_2D
	Texture2D = C.GL_TEXTURE_2D
	// Blend GL_BLEND
	Blend            = C.GL_BLEND
	SrcAlpha         = C.GL_SRC_ALPHA
	One              = C.GL_ONE
	OneMinusSrcAlpha = C.GL_ONE_MINUS_SRC_ALPHA
)

func boolToGL(val bool) C.uchar {
	if val {
		return C.GL_TRUE
	}
	return C.GL_FALSE
}

// Enable glEnable
func Enable(cap uint) {
	C.glEnable(C.uint(cap))
}

// Disable glDisable
func Disable(cap uint) {
	C.glDisable(C.uint(cap))
}

// LineWidth glLineWidth
func LineWidth(width float32) {
	C.glLineWidth(C.float(width))
}

// DepthMask glDepthMask
func DepthMask(flag bool) {
	C.glDepthMask(boolToGL(flag))
}

// BlendFunc glBlendFunc
func BlendFunc(sfactor, dfactor uint) {
	C.glBlendFunc(C.uint(sfactor), C.uint(dfactor))
}

// Color4f glColor4f
func Color4f(r, g, b, a float32) {
	C.glColor4f(C.float(r), C.float(g), C.float(b), C.float(a))
}
