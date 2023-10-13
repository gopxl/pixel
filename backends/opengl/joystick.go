package opengl

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gopxl/pixel/v2"
)

var joystickMapping = map[glfw.Joystick]pixel.Joystick{
	glfw.Joystick1:  pixel.Joystick1,
	glfw.Joystick2:  pixel.Joystick2,
	glfw.Joystick3:  pixel.Joystick3,
	glfw.Joystick4:  pixel.Joystick4,
	glfw.Joystick5:  pixel.Joystick5,
	glfw.Joystick6:  pixel.Joystick6,
	glfw.Joystick7:  pixel.Joystick7,
	glfw.Joystick8:  pixel.Joystick8,
	glfw.Joystick9:  pixel.Joystick9,
	glfw.Joystick10: pixel.Joystick10,
	glfw.Joystick11: pixel.Joystick11,
	glfw.Joystick12: pixel.Joystick12,
	glfw.Joystick13: pixel.Joystick13,
	glfw.Joystick14: pixel.Joystick14,
	glfw.Joystick15: pixel.Joystick15,
	glfw.Joystick16: pixel.Joystick16,
}

// Not currently used because Gamepad Axis/Button input works a bit different than others
var _ = map[glfw.GamepadAxis]pixel.GamepadAxis{
	glfw.AxisLeftX:        pixel.AxisLeftX,
	glfw.AxisLeftY:        pixel.AxisLeftY,
	glfw.AxisRightX:       pixel.AxisRightX,
	glfw.AxisRightY:       pixel.AxisRightY,
	glfw.AxisLeftTrigger:  pixel.AxisLeftTrigger,
	glfw.AxisRightTrigger: pixel.AxisRightTrigger,
}

var _ = map[glfw.GamepadButton]pixel.GamepadButton{
	glfw.ButtonA:           pixel.GamepadA,
	glfw.ButtonB:           pixel.GamepadB,
	glfw.ButtonX:           pixel.GamepadX,
	glfw.ButtonY:           pixel.GamepadY,
	glfw.ButtonLeftBumper:  pixel.GamepadLeftBumper,
	glfw.ButtonRightBumper: pixel.GamepadRightBumper,
	glfw.ButtonBack:        pixel.GamepadBack,
	glfw.ButtonStart:       pixel.GamepadStart,
	glfw.ButtonGuide:       pixel.GamepadGuide,
	glfw.ButtonLeftThumb:   pixel.GamepadLeftThumb,
	glfw.ButtonRightThumb:  pixel.GamepadRightThumb,
	glfw.ButtonDpadUp:      pixel.GamepadDpadUp,
	glfw.ButtonDpadRight:   pixel.GamepadDpadRight,
	glfw.ButtonDpadDown:    pixel.GamepadDpadDown,
	glfw.ButtonDpadLeft:    pixel.GamepadDpadLeft,
}

// JoystickPresent returns if the joystick is currently connected.
//
// This API is experimental.
func (w *Window) JoystickPresent(js pixel.Joystick) bool {
	return w.currJoy.connected[js]
}

// JoystickName returns the name of the joystick. A disconnected joystick will return an
// empty string.
//
// This API is experimental.
func (w *Window) JoystickName(js pixel.Joystick) string {
	return w.currJoy.name[js]
}

// JoystickButtonCount returns the number of buttons a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(js pixel.Joystick) int {
	return len(w.currJoy.buttons[js])
}

// JoystickAxisCount returns the number of axes a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(js pixel.Joystick) int {
	return len(w.currJoy.axis[js])
}

// JoystickPressed returns whether the joystick Button is currently pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickPressed(js pixel.Joystick, button pixel.GamepadButton) bool {
	return w.currJoy.getButton(js, int(button))
}

// JoystickJustPressed returns whether the joystick Button has just been pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(js pixel.Joystick, button pixel.GamepadButton) bool {
	return w.currJoy.getButton(js, int(button)) && !w.prevJoy.getButton(js, int(button))
}

// JoystickJustReleased returns whether the joystick Button has just been released up.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(js pixel.Joystick, button pixel.GamepadButton) bool {
	return !w.currJoy.getButton(js, int(button)) && w.prevJoy.getButton(js, int(button))
}

// JoystickAxis returns the value of a joystick axis at the last call to Window.Update.
// If the axis index is out of range, this will return 0.
//
// This API is experimental.
func (w *Window) JoystickAxis(js pixel.Joystick, axis pixel.GamepadAxis) float64 {
	return w.currJoy.getAxis(js, int(axis))
}

// Used internally during Window.UpdateInput to update the state of the joysticks.
func (w *Window) updateJoystickInput() {
	for joystick := glfw.Joystick1; joystick <= glfw.JoystickLast; joystick++ {
		js, ok := joystickMapping[joystick]
		if !ok {
			return
		}
		// Determine and store if the joystick was connected
		joystickPresent := joystick.Present()
		w.tempJoy.connected[js] = joystickPresent

		if joystickPresent {
			if joystick.IsGamepad() {
				gamepadInputs := joystick.GetGamepadState()

				w.tempJoy.buttons[js] = gamepadInputs.Buttons[:]
				w.tempJoy.axis[js] = gamepadInputs.Axes[:]
			} else {
				w.tempJoy.buttons[js] = joystick.GetButtons()
				w.tempJoy.axis[js] = joystick.GetAxes()
			}

			if !w.currJoy.connected[js] {
				// The joystick was recently connected, we get the name
				w.tempJoy.name[js] = joystick.GetName()
			} else {
				// Use the name from the previous one
				w.tempJoy.name[js] = w.currJoy.name[js]
			}
		} else {
			w.tempJoy.buttons[js] = []glfw.Action{}
			w.tempJoy.axis[js] = []float32{}
			w.tempJoy.name[js] = ""
		}
	}

	w.prevJoy = w.currJoy
	w.currJoy = w.tempJoy
}

type joystickState struct {
	connected [pixel.JoystickLast + 1]bool
	name      [pixel.JoystickLast + 1]string
	buttons   [pixel.JoystickLast + 1][]glfw.Action
	axis      [pixel.JoystickLast + 1][]float32
}

// Returns if a button on a joystick is down, returning false if the button or joystick is invalid.
func (js *joystickState) getButton(joystick pixel.Joystick, button int) bool {
	// Check that the joystick and button is valid, return false by default
	if js.buttons[joystick] == nil || button >= len(js.buttons[joystick]) || button < 0 {
		return false
	}
	return js.buttons[joystick][byte(button)] == glfw.Press
}

// Returns the value of a joystick axis, returning 0 if the button or joystick is invalid.
func (js *joystickState) getAxis(joystick pixel.Joystick, axis int) float64 {
	// Check that the joystick and axis is valid, return 0 by default.
	if js.axis[joystick] == nil || axis >= len(js.axis[joystick]) || axis < 0 {
		return 0
	}
	return float64(js.axis[joystick][axis])
}
