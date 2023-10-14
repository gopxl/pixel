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
func (w *Window) JoystickPresent(js glfw.Joystick) bool {
	return w.currJoy[js].Connected()
}

// JoystickName returns the name of the joystick. A disconnected joystick will return an
// empty string.
//
// This API is experimental.
func (w *Window) JoystickName(js glfw.Joystick) string {
	return w.currJoy[js].Name()
}

// JoystickButtonCount returns the number of buttons a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(js glfw.Joystick) int {
	return w.currJoy[js].NumButtons()
}

// JoystickAxisCount returns the number of axes a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(js glfw.Joystick) int {
	return w.currJoy[js].NumAxes()
}

// JoystickPressed returns whether the joystick Button is currently pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickPressed(js glfw.Joystick, button glfw.GamepadButton) bool {
	return w.currJoy[js].Button(button) == glfw.Press
}

// JoystickJustPressed returns whether the joystick Button has just been pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(js glfw.Joystick, button glfw.GamepadButton) bool {
	return w.currJoy[js].Button(button) == glfw.Press && w.prevJoy[js].Button(button) != glfw.Press
}

// JoystickJustReleased returns whether the joystick Button has just been released up.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(js glfw.Joystick, button glfw.GamepadButton) bool {
	return w.currJoy[js].Button(button) != glfw.Press && w.prevJoy[js].Button(button) == glfw.Press
}

// JoystickAxis returns the value of a joystick axis at the last call to Window.Update.
// If the axis index is out of range, this will return 0.
//
// This API is experimental.
func (w *Window) JoystickAxis(js glfw.Joystick, axis glfw.GamepadAxis) float64 {
	return float64(w.currJoy[js].Axis(axis))
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
		w.tempJoy[js].SetConnected(joystickPresent)

		if joystickPresent {
			if joystick.IsGamepad() {
				gamepadInputs := joystick.GetGamepadState()

				w.tempJoy[js].SetButtons(gamepadInputs.Buttons[:])
				w.tempJoy[js].SetAxes(gamepadInputs.Axes[:])
			} else {
				w.tempJoy[js].SetButtons(joystick.GetButtons())
				w.tempJoy[js].SetAxes(joystick.GetAxes())
			}

			if !w.currJoy[js].Connected() {
				// The joystick was recently connected, we get the name
				w.tempJoy[js].SetName(joystick.GetName())
			} else {
				// Use the name from the previous one
				w.tempJoy[js].SetName(w.currJoy[js].Name())
			}
		} else {
			w.tempJoy[js].SetButtons([]glfw.Action{})
			w.tempJoy[js].SetAxes([]float32{})
			w.tempJoy[js].SetName("")
		}
	}

	w.prevJoy = w.currJoy
	w.currJoy = w.tempJoy
}
