# gl.4413a [![Build Status](https://github.com/go-gl/gl/actions/workflows/main.yml/badge.svg)](https://github.com/go-gl/gl/actions/workflows/main.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/go-gl/gl.svg)](https://pkg.go.dev/github.com/go-gl/gl)

This repository holds Go bindings to various OpenGL versions. They are auto-generated using [Glow](https://github.com/go-gl/glow).

Features:
- Go functions that mirror the C specification using Go types.
- Support for multiple OpenGL APIs (GL/GLES/EGL/WGL/GLX/EGL), versions, and profiles.
- Support for extensions (including debug callbacks).

Requirements:
- A cgo compiler (typically gcc).
- On Ubuntu/Debian-based systems, the `libgl1-mesa-dev` package.

## Usage

Use `go get -u` to download and install the prebuilt packages. The prebuilt packages support OpenGL versions 2.1, 3.1, 3.2, 3.3, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6 across both the core and compatibility profiles and include all extensions. Pick whichever one(s) you need:

    go get -u github.com/go-gl/gl/v{3.2,3.3,4.1,4.2,4.3,4.4,4.5,4.6}-{core,compatibility}/gl
    go get -u github.com/go-gl/gl/v3.1/gles2
    go get -u github.com/go-gl/gl/v2.1/gl

Once the bindings are installed you can use them with the appropriate import statements.

```Go
package main

import "github.com/go-gl/gl/v3.3-core/gl"

func main() {
	window := ... // Open a window.
	window.MakeContextCurrent()

	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
}
```

The `gl` package contains the OpenGL functions and enumeration values for the imported version. It also contains helper functions for working with the API. Of note is `gl.Ptr` which takes a Go array or slice or pointer and returns a corresponding `uintptr` to use with functions expecting data pointers. Also of note is `gl.Str` which takes a null-terminated Go string and returns a corresponding `*int8` to use with functions expecting character pointers.

A note about threading and goroutines. The bindings do not expose a mechanism to make an OpenGL context current on a different thread so you must restrict your usage to the thread on which you called `gl.Init()`. To do so you should use [LockOSThread](https://code.google.com/p/go-wiki/wiki/LockOSThread).

## Examples

Examples illustrating how to use the bindings are available in the [example](https://github.com/go-gl/example) repo. There are examples for [OpenGL 4.1 core](https://github.com/go-gl/example/tree/master/gl41core-cube) and [OpenGL 2.1](https://github.com/go-gl/example/tree/master/gl21-cube).

## Function Loading

The `procaddr` package contains platform-specific functions for [loading OpenGL functions](https://www.opengl.org/wiki/Load_OpenGL_Functions). Calling `gl.Init()` uses the `auto` subpackage to automatically select an appropriate implementation based on the build environment. If you want to select a specific implementation you can use the `noauto` build tag and the `gl.InitWithProcAddrFunc` initialization function.

## Go >=1.14 and `checkptr`

In version 1.14 of Go, the race detector added `checkptr` instrumentation. This compilation option ensures that programs follow `unsafe.Pointer` safety rules. See here for details: https://golang.org/doc/go1.14#compiler.

If enabled, there is a high chance that it will cause program termination when calling specific OpenGL functions, with a message like this:

```
fatal error: checkptr: pointer arithmetic computed bad pointer value
```

The reported call stack will point to a function like `gl.VertexAttribPointer()` that receives an `unsafe.Pointer` as parameter.
In case such function requires an "offset" passed in as a pointer (in the low-level API), a different signature needs to be used in order to satisfy the detector.

For this purpose `glow` generates "override" functions which have a different signature, taking `uintptr` instead. These functions have the suffix `WithOffset` (or similar) in their name.
For the previous example, it would be `gl.VertexAttribPointerWithOffset()`.

Not all such functions have an appropriate override! In case you stumble over such an error, and the override is missing, you have the following options:
* Disable the detector by building your program with `-gcflags=all=-d=checkptr=0`
* Report the missing function(s) as [issue for `glow`](https://github.com/go-gl/glow/issues)
* Possibly even create a pull-request for `glow` with the missing override yourself, and re-generate the `gl` bindings.

## Generating

These gl bindings are generated using the [Glow](https://github.com/go-gl/glow) generator. Only developers of this repository need to do this step.

It is required to have `glow` source in a sibling directory to `go-gl/gl` since relative paths are used for generation (see `generate.go`).
For non-module-aware cases, this means `glow` needs to be in the same Go workspace as `go-gl/gl`.
For module-aware cases, `go-gl/glow` needs to be checked out parallel to `go-gl/gl`. 

Perform generation with the following:

```bash
cd path/to/go-gl/gl
go generate -tags=gen .
```

More information about these bindings can be found in the [Glow repository](https://github.com/go-gl/glow).
