//go:build js && wasm

package opengl

import "github.com/gopxl/mainthread/v2"

// Run invokes the supplied function on the JS event-loop thread. The WASM
// mainthread shim calls the function inline, so this is a thin pass-through
// that keeps the call-site compatible with the desktop backend.
func Run(run func()) {
	mainthread.Run(run)
}
