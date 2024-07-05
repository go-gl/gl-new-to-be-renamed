package gl_test

import (
	"errors"
	"runtime"
	"testing"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func testIntegers(t *testing.T) {
	// See https://registry.khronos.org/OpenGL-Refpages/gl4/html/glGet.xhtml
	var data int32
	gl.GetIntegerv(gl.MAJOR_VERSION, &data)
	if data != 4 {
		// OpenGL 5.0 released with raytracing...?
		t.Error("invalid GL_MAJOR_VERSION:", data)
	}
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &data)
	if data < 1024 {
		// Guaranteed by spec
		t.Error("invalid GL_MAX_TEXTURE_SIZE:", data)
	}

	if err := gl.GetError(); err != gl.NO_ERROR {
		t.Error("glGetIntegerv():", err)
	}
}
func testStrings(t *testing.T) {
	// See https://registry.khronos.org/OpenGL-Refpages/gl4/html/glGetString.xhtml
	gl.GetString(gl.VENDOR)
	gl.GetString(gl.RENDERER)
	gl.GetString(gl.VERSION)
	gl.GetString(gl.SHADING_LANGUAGE_VERSION)
	if err := gl.GetError(); err != gl.NO_ERROR {
		t.Error("glGetString():", err)
	}

	gl.GetString(gl.MAX_TEXTURE_SIZE)
	if err := gl.GetError(); err != gl.INVALID_ENUM {
		t.Error("glGetString() failed to return GL_INVALID_ENUM:", err)
	}
}
func testTextures(t *testing.T) {
	var texture uint32
	gl.GenTextures(1, &texture)
	if texture == 0 {
		t.Error("glGenTextures() returned zero")
	}

	// Textures must be bound before glIsTexture will recognize them.
	// See https://registry.khronos.org/OpenGL-Refpages/gl4/html/glIsTexture.xhtml
	gl.BindTexture(gl.TEXTURE_2D, texture)
	if !gl.IsTexture(texture) {
		t.Error("glIsTexture() failed to recognize a texture returned by glGenTextures()")
	}

	gl.DeleteTextures(1, &texture)
	if gl.IsTexture(texture) {
		t.Error("glDeleteTextures() did not delete texture")
	}

	if err := gl.GetError(); err != gl.NO_ERROR {
		t.Error("texture error:", err)
	}

	gl.GenTextures(-1, &texture)
	if err := gl.GetError(); err != gl.INVALID_VALUE {
		t.Error("glGenTextures() failed to return GL_INVALID_VALUE:", err)
	}
}
func testShader(t *testing.T, src string) error {
	csrc, free := gl.Strs(src + "\x00")
	defer free()

	shader := gl.CreateShader(gl.VERTEX_SHADER)
	if shader == 0 {
		t.Error("glCreateShader() returned zero")
	}
	defer gl.DeleteShader(shader)

	gl.ShaderSource(shader, 1, csrc, nil)
	gl.CompileShader(shader)
	if !gl.IsShader(shader) {
		t.Error("glIsShader() failed to recognize a shader returned by glCreateShader()")
	}

	if err := gl.GetError(); err != gl.NO_ERROR {
		t.Error("shader error:", err)
	}

	var data int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &data)
	if data == gl.TRUE {
		return nil
	}

	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &data)
	infoLog := make([]byte, data+1)
	gl.GetShaderInfoLog(shader, data, nil, &infoLog[0])
	return errors.New(src + "\n" + string(infoLog))
}

func TestBasic(t *testing.T) {
	// Each test runs in its own goroutine, so we need to lock OS threads here
	// rather than in the init() function.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := glfw.Init(); err != nil {
		t.Fatal("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// Needed for OS X
	// https://www.glfw.org/faq.html#41---how-do-i-create-an-opengl-30-context
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "Test", nil, nil)
	if err != nil {
		t.Fatal("failed to create glfw window:", err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		t.Fatal("failed to initialize opengl:", err)
	}

	testIntegers(t)
	testStrings(t)
	testTextures(t)

	err = testShader(t, `
		#version 410
		void main() {
			gl_Position = vec4(0, 0, 0, 1);
		}
	`)
	if err != nil {
		t.Error("unexpected compile error:", err)
	}

	err = testShader(t, `
		#version 410
		void main() {
			gl_Unlucky = vec13(13);
		}
	`)
	if err == nil {
		t.Error("unexpected successful compilation of invalid shader")
	} else {
		t.Log(err)
	}
}
