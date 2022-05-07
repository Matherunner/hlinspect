package game

// Very rudimentary and adhoc Go bindings to OpenGL

/*
#cgo 386 LDFLAGS: -lopengl32

#include <GL/gl.h>
*/
import "C"

type GLEnum uint32

const (
	GLTrue             GLEnum = C.GL_TRUE
	GLFalse            GLEnum = C.GL_FALSE
	GLTexture2D        GLEnum = C.GL_TEXTURE_2D
	GLBlend            GLEnum = C.GL_BLEND
	GLSrcAlpha         GLEnum = C.GL_SRC_ALPHA
	GLOne              GLEnum = C.GL_ONE
	GLOneMinusSrcAlpha GLEnum = C.GL_ONE_MINUS_SRC_ALPHA
	GLColorBufferBit   GLEnum = C.GL_COLOR_BUFFER_BIT
)

type GL struct {
}

func NewGL() *GL {
	return &GL{}
}

func (gl *GL) boolToGL(val bool) C.uchar {
	if val {
		return C.GL_TRUE
	}
	return C.GL_FALSE
}

func (gl *GL) Enable(cap GLEnum) {
	C.glEnable(C.uint(cap))
}

func (gl *GL) Disable(cap GLEnum) {
	C.glDisable(C.uint(cap))
}

func (gl *GL) LineWidth(width float32) {
	C.glLineWidth(C.float(width))
}

func (gl *GL) DepthMask(flag bool) {
	C.glDepthMask(gl.boolToGL(flag))
}

func (gl *GL) BlendFunc(sfactor, dfactor GLEnum) {
	C.glBlendFunc(C.uint(sfactor), C.uint(dfactor))
}

func (gl *GL) Color4f(r, g, b, a float32) {
	C.glColor4f(C.float(r), C.float(g), C.float(b), C.float(a))
}

func (gl *GL) Clear(mask GLEnum) {
	C.glClear(C.uint(mask))
}

func (gl *GL) ClearColor(r, g, b, a float32) {
	C.glClearColor(C.float(r), C.float(g), C.float(b), C.float(a))
}
