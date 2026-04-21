//go:build js && wasm

package opengl

import (
	"time"

	"github.com/gopxl/pixel/v2"
)

// Input state is tracked via the same InputHandler the desktop backend uses;
// DOM event listeners installed by Window.initInput feed it through
// ButtonEvent / CharEvent / MouseMoveEvent / MouseScrollEvent.

func (w *Window) Pressed(button pixel.Button) bool      { return w.input.Curr.Buttons[button] }
func (w *Window) JustPressed(button pixel.Button) bool  { return w.input.PressEvents[button] }
func (w *Window) JustReleased(button pixel.Button) bool { return w.input.ReleaseEvents[button] }
func (w *Window) Repeated(button pixel.Button) bool     { return w.input.Curr.Repeat[button] }

func (w *Window) MousePosition() pixel.Vec         { return w.input.Curr.Mouse }
func (w *Window) MousePreviousPosition() pixel.Vec { return w.input.Prev.Mouse }
func (w *Window) SetMousePosition(v pixel.Vec)     { w.input.SetMousePosition(v) }
func (w *Window) MouseInsideWindow() bool          { return w.input.MouseInsideWindow }
func (w *Window) MouseScroll() pixel.Vec           { return w.input.Curr.Scroll }
func (w *Window) MousePreviousScroll() pixel.Vec   { return w.input.Prev.Scroll }
func (w *Window) Typed() string                    { return w.input.Curr.Typed }

func (w *Window) SetButtonCallback(cb func(win *Window, button pixel.Button, action pixel.Action)) {
	w.buttonCallback = cb
}
func (w *Window) SetCharCallback(cb func(win *Window, r rune)) { w.charCallback = cb }
func (w *Window) SetMouseEnteredCallback(cb func(win *Window, entered bool)) {
	w.mouseEnteredCallback = cb
}
func (w *Window) SetMouseMovedCallback(cb func(win *Window, pos pixel.Vec)) {
	w.mouseMovedCallback = cb
}
func (w *Window) SetScrollCallback(cb func(win *Window, scroll pixel.Vec)) {
	w.scrollCallback = cb
}

// UpdateInput commits the pending input events for the next frame and polls
// the connected gamepads via the HTML5 Gamepad API.
func (w *Window) UpdateInput() {
	w.input.Update()
	w.updateJoystickInput()
}

// UpdateInputWait commits pending events. The timeout argument is accepted for
// API parity with the desktop backend but ignored under WASM — the browser
// event loop delivers events asynchronously.
func (w *Window) UpdateInputWait(timeout time.Duration) {
	w.input.Update()
	w.updateJoystickInput()
}
