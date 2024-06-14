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

var gamepadButtonMapping = map[glfw.GamepadButton]pixel.GamepadButton{
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
	return w.currJoy.Connected[js]
}

// JoystickName returns the name of the joystick. A disconnected joystick will return an
// empty string.
//
// This API is experimental.
func (w *Window) JoystickName(js pixel.Joystick) string {
	return w.currJoy.Name[js]
}

// JoystickButtonCount returns the number of buttons a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(js pixel.Joystick) int {
	return len(w.currJoy.Buttons[js])
}

// JoystickAxisCount returns the number of axes a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(js pixel.Joystick) int {
	return len(w.currJoy.Axis[js])
}

// JoystickPressed returns whether the joystick Button is currently pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickPressed(js pixel.Joystick, button pixel.GamepadButton) bool {
	return w.currJoy.GetButton(js, button)
}

// JoystickJustPressed returns whether the joystick Button has just been pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(js pixel.Joystick, button pixel.GamepadButton) bool {
	return w.currJoy.GetButton(js, button) && !w.prevJoy.GetButton(js, button)
}

// JoystickJustReleased returns whether the joystick Button has just been released up.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(js pixel.Joystick, button pixel.GamepadButton) bool {
	return !w.currJoy.GetButton(js, button) && w.prevJoy.GetButton(js, button)
}

// JoystickAxis returns the value of a joystick axis at the last call to Window.Update.
// If the axis index is out of range, this will return 0.
//
// This API is experimental.
func (w *Window) JoystickAxis(js pixel.Joystick, axis pixel.GamepadAxis) float64 {
	return w.currJoy.GetAxis(js, axis)
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
		w.tempJoy.Connected[js] = joystickPresent

		if joystickPresent {
			if joystick.IsGamepad() {
				gamepadInputs := joystick.GetGamepadState()

				w.tempJoy.Buttons[js] = convertGamepadButtons(gamepadInputs.Buttons)
				w.tempJoy.Axis[js] = gamepadInputs.Axes[:]
			} else {
				w.tempJoy.Buttons[js] = convertJoystickButtons(joystick.GetButtons())
				w.tempJoy.Axis[js] = joystick.GetAxes()
			}

			if !w.currJoy.Connected[js] {
				// The joystick was recently connected, we get the name
				w.tempJoy.Name[js] = joystick.GetName()
			} else {
				// Use the name from the previous one
				w.tempJoy.Name[js] = w.currJoy.Name[js]
			}
		} else {
			w.tempJoy.Buttons[js] = []pixel.Action{}
			w.tempJoy.Axis[js] = []float32{}
			w.tempJoy.Name[js] = ""
		}
	}

	w.prevJoy = w.currJoy
	w.currJoy = w.tempJoy
}

// Convert buttons from a GLFW gamepad mapping to pixel format
func convertGamepadButtons(buttons [glfw.ButtonLast + 1]glfw.Action) []pixel.Action {
	pixelButtons := make([]pixel.Action, pixel.NumGamepadButtons)
	for i, a := range buttons {
		var action pixel.Action
		var button pixel.GamepadButton
		var ok bool
		if action, ok = actionMapping[a]; !ok {
			// Unknown action
			continue
		}
		if button, ok = gamepadButtonMapping[glfw.GamepadButton(i)]; !ok {
			// Unknown gamepad button
			continue
		}
		pixelButtons[button] = action
	}
	return pixelButtons
}

// Convert buttons of unknown length and arrangement to pixel format
// Used when a joystick has an unknown mapping in GLFW
func convertJoystickButtons(buttons []glfw.Action) []pixel.Action {
	pixelButtons := make([]pixel.Action, len(buttons))
	for i, a := range buttons {
		var action pixel.Action
		var ok bool
		if action, ok = actionMapping[a]; !ok {
			// Unknown action
			continue
		}
		pixelButtons[pixel.GamepadButton(i)] = action
	}
	return pixelButtons
}
