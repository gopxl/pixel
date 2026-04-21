//go:build js && wasm

package opengl

import (
	"syscall/js"

	"github.com/gopxl/pixel/v2"
)

// Browser "standard" gamepad button indices (see
// https://www.w3.org/TR/gamepad/#remapping).
const (
	webBtnA            = 0
	webBtnB            = 1
	webBtnX            = 2
	webBtnY            = 3
	webBtnLeftBumper   = 4
	webBtnRightBumper  = 5
	webBtnLeftTrigger  = 6
	webBtnRightTrigger = 7
	webBtnBack         = 8
	webBtnStart        = 9
	webBtnLeftThumb    = 10
	webBtnRightThumb   = 11
	webBtnDpadUp       = 12
	webBtnDpadDown     = 13
	webBtnDpadLeft     = 14
	webBtnDpadRight    = 15
	webBtnGuide        = 16
)

// webButtonToPixel maps the browser's standard-layout button indices onto
// pixel's GamepadButton values. Triggers (6, 7) are intentionally omitted
// here and surface as AxisLeftTrigger / AxisRightTrigger so the API matches
// the desktop GLFW backend.
var webButtonToPixel = map[int]pixel.GamepadButton{
	webBtnA:           pixel.GamepadA,
	webBtnB:           pixel.GamepadB,
	webBtnX:           pixel.GamepadX,
	webBtnY:           pixel.GamepadY,
	webBtnLeftBumper:  pixel.GamepadLeftBumper,
	webBtnRightBumper: pixel.GamepadRightBumper,
	webBtnBack:        pixel.GamepadBack,
	webBtnStart:       pixel.GamepadStart,
	webBtnGuide:       pixel.GamepadGuide,
	webBtnLeftThumb:   pixel.GamepadLeftThumb,
	webBtnRightThumb:  pixel.GamepadRightThumb,
	webBtnDpadUp:      pixel.GamepadDpadUp,
	webBtnDpadRight:   pixel.GamepadDpadRight,
	webBtnDpadDown:    pixel.GamepadDpadDown,
	webBtnDpadLeft:    pixel.GamepadDpadLeft,
}

// JoystickPresent reports whether a gamepad is connected in the given slot.
//
// This API is experimental.
func (w *Window) JoystickPresent(j pixel.Joystick) bool {
	return w.currJoy.Connected[j]
}

// JoystickName returns the navigator-supplied id for the gamepad in the given
// slot, or an empty string if no gamepad is present.
//
// This API is experimental.
func (w *Window) JoystickName(j pixel.Joystick) string {
	return w.currJoy.Name[j]
}

// JoystickButtonCount returns the number of buttons exposed by the connected
// gamepad. Returns 0 for disconnected slots.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(j pixel.Joystick) int {
	return len(w.currJoy.Buttons[j])
}

// JoystickAxisCount returns the number of axes exposed by the connected
// gamepad. Returns 0 for disconnected slots.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(j pixel.Joystick) int {
	return len(w.currJoy.Axis[j])
}

// JoystickPressed reports whether the given button is currently held.
//
// This API is experimental.
func (w *Window) JoystickPressed(j pixel.Joystick, button pixel.GamepadButton) bool {
	return w.currJoy.GetButton(j, button)
}

// JoystickJustPressed reports whether the given button transitioned to
// pressed since the last UpdateInput.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(j pixel.Joystick, button pixel.GamepadButton) bool {
	return w.currJoy.GetButton(j, button) && !w.prevJoy.GetButton(j, button)
}

// JoystickJustReleased reports whether the given button transitioned to
// released since the last UpdateInput.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(j pixel.Joystick, button pixel.GamepadButton) bool {
	return !w.currJoy.GetButton(j, button) && w.prevJoy.GetButton(j, button)
}

// JoystickAxis returns the current value of the given axis, in [-1, 1] for
// sticks and [0, 1] for triggers.
//
// This API is experimental.
func (w *Window) JoystickAxis(j pixel.Joystick, axis pixel.GamepadAxis) float64 {
	return w.currJoy.GetAxis(j, axis)
}

// updateJoystickInput polls navigator.getGamepads() and rotates the joystick
// state snapshots. Called once per frame from UpdateInput.
//
// The Gamepad API is poll-based: the browser updates snapshots asynchronously
// and the page observes state by re-reading the list each frame.
func (w *Window) updateJoystickInput() {
	pads := browserGamepads()
	for slot := 0; slot < pixel.NumJoysticks; slot++ {
		joy := pixel.Joystick(slot)
		if !pads.Truthy() || slot >= pads.Length() {
			w.clearJoySlot(joy)
			continue
		}
		pad := pads.Index(slot)
		if !pad.Truthy() || !pad.Get("connected").Bool() {
			w.clearJoySlot(joy)
			continue
		}
		w.tempJoy.Connected[joy] = true
		w.tempJoy.Name[joy] = pad.Get("id").String()
		w.tempJoy.Buttons[joy], w.tempJoy.Axis[joy] = readPadState(pad)
	}
	w.prevJoy = w.currJoy
	w.currJoy = w.tempJoy
}

func (w *Window) clearJoySlot(j pixel.Joystick) {
	w.tempJoy.Connected[j] = false
	w.tempJoy.Buttons[j] = nil
	w.tempJoy.Axis[j] = nil
	w.tempJoy.Name[j] = ""
}

// browserGamepads calls navigator.getGamepads() and returns the resulting
// array. Returns js.Undefined() when the host lacks the Gamepad API (older
// browsers, non-secure contexts, headless environments).
func browserGamepads() js.Value {
	nav := js.Global().Get("navigator")
	if !nav.Truthy() {
		return js.Undefined()
	}
	fn := nav.Get("getGamepads")
	if !fn.Truthy() {
		return js.Undefined()
	}
	return nav.Call("getGamepads")
}

// readPadState converts a browser Gamepad object into pixel's
// (buttons, axes) representation. Standard-layout pads are remapped so the
// button and axis order matches the desktop GLFW backend; non-standard pads
// pass through untouched so applications can still address raw indices.
func readPadState(pad js.Value) ([]pixel.Action, []float32) {
	mapping := pad.Get("mapping").String()
	btnsJS := pad.Get("buttons")
	axesJS := pad.Get("axes")
	btnN := btnsJS.Length()
	axN := axesJS.Length()

	if mapping == "standard" {
		buttons := make([]pixel.Action, pixel.NumGamepadButtons)
		for i := 0; i < btnN; i++ {
			target, ok := webButtonToPixel[i]
			if !ok {
				continue
			}
			if btnsJS.Index(i).Get("pressed").Bool() {
				buttons[target] = pixel.Press
			} else {
				buttons[target] = pixel.Release
			}
		}
		// Standard axes are [LX, LY, RX, RY]. Pixel also exposes the two
		// triggers as axes, but on the web they are reported as analog
		// buttons (indices 6 and 7) — promote them so downstream code sees
		// the same axis layout as on desktop.
		axes := make([]float32, pixel.NumAxes)
		for i := 0; i < axN && i < 4; i++ {
			axes[i] = float32(axesJS.Index(i).Float())
		}
		if btnN > webBtnLeftTrigger {
			axes[pixel.AxisLeftTrigger] = float32(btnsJS.Index(webBtnLeftTrigger).Get("value").Float())
		}
		if btnN > webBtnRightTrigger {
			axes[pixel.AxisRightTrigger] = float32(btnsJS.Index(webBtnRightTrigger).Get("value").Float())
		}
		return buttons, axes
	}

	buttons := make([]pixel.Action, btnN)
	for i := 0; i < btnN; i++ {
		if btnsJS.Index(i).Get("pressed").Bool() {
			buttons[i] = pixel.Press
		} else {
			buttons[i] = pixel.Release
		}
	}
	axes := make([]float32, axN)
	for i := 0; i < axN; i++ {
		axes[i] = float32(axesJS.Index(i).Float())
	}
	return buttons, axes
}
