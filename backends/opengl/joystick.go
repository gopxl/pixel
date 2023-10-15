package opengl

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gopxl/pixel/v2"
)

var joystickMapping = map[pixel.Joystick]glfw.Joystick{
	pixel.Joystick1:  glfw.Joystick1,
	pixel.Joystick2:  glfw.Joystick2,
	pixel.Joystick3:  glfw.Joystick3,
	pixel.Joystick4:  glfw.Joystick4,
	pixel.Joystick5:  glfw.Joystick5,
	pixel.Joystick6:  glfw.Joystick6,
	pixel.Joystick7:  glfw.Joystick7,
	pixel.Joystick8:  glfw.Joystick8,
	pixel.Joystick9:  glfw.Joystick9,
	pixel.Joystick10: glfw.Joystick10,
	pixel.Joystick11: glfw.Joystick11,
	pixel.Joystick12: glfw.Joystick12,
	pixel.Joystick13: glfw.Joystick13,
	pixel.Joystick14: glfw.Joystick14,
	pixel.Joystick15: glfw.Joystick15,
	pixel.Joystick16: glfw.Joystick16,
}

// Not currently used because Gamepad Axis/Button input works a bit different than others
var gamepadAxisMapping = map[pixel.GamepadAxis]glfw.GamepadAxis{
	pixel.AxisLeftX:        glfw.AxisLeftX,
	pixel.AxisLeftY:        glfw.AxisLeftY,
	pixel.AxisRightX:       glfw.AxisRightX,
	pixel.AxisRightY:       glfw.AxisRightY,
	pixel.AxisLeftTrigger:  glfw.AxisLeftTrigger,
	pixel.AxisRightTrigger: glfw.AxisRightTrigger,
}

var gamepadButtonMapping = map[pixel.GamepadButton]glfw.GamepadButton{
	pixel.GamepadA:           glfw.ButtonA,
	pixel.GamepadB:           glfw.ButtonB,
	pixel.GamepadX:           glfw.ButtonX,
	pixel.GamepadY:           glfw.ButtonY,
	pixel.GamepadLeftBumper:  glfw.ButtonLeftBumper,
	pixel.GamepadRightBumper: glfw.ButtonRightBumper,
	pixel.GamepadBack:        glfw.ButtonBack,
	pixel.GamepadStart:       glfw.ButtonStart,
	pixel.GamepadGuide:       glfw.ButtonGuide,
	pixel.GamepadLeftThumb:   glfw.ButtonLeftThumb,
	pixel.GamepadRightThumb:  glfw.ButtonRightThumb,
	pixel.GamepadDpadUp:      glfw.ButtonDpadUp,
	pixel.GamepadDpadRight:   glfw.ButtonDpadRight,
	pixel.GamepadDpadDown:    glfw.ButtonDpadDown,
	pixel.GamepadDpadLeft:    glfw.ButtonDpadLeft,
}

// JoystickPresent returns if the joystick is currently connected.
//
// This API is experimental.
func (w *Window) JoystickPresent(js pixel.Joystick) bool {
	return w.currJoy[js].Connected()
}

// JoystickName returns the name of the joystick. A disconnected joystick will return an
// empty string.
//
// This API is experimental.
func (w *Window) JoystickName(js pixel.Joystick) string {
	return w.currJoy[js].Name()
}

// JoystickButtonCount returns the number of buttons a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(js pixel.Joystick) int {
	return w.currJoy[js].NumButtons()
}

// JoystickAxisCount returns the number of axes a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(js pixel.Joystick) int {
	return w.currJoy[js].NumAxes()
}

// JoystickPressed returns whether the joystick Button is currently pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickPressed(js pixel.Joystick, button pixel.GamepadButton) bool {
	if b, ok := gamepadButtonMapping[button]; ok {
		return w.currJoy[js].Button(b) == glfw.Press
	}
	return false
}

// JoystickJustPressed returns whether the joystick Button has just been pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(js pixel.Joystick, button pixel.GamepadButton) bool {
	if b, ok := gamepadButtonMapping[button]; ok {
		return w.currJoy[js].Button(b) == glfw.Press && w.prevJoy[js].Button(b) != glfw.Press
	}
	return false
}

// JoystickJustReleased returns whether the joystick Button has just been released up.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(js pixel.Joystick, button pixel.GamepadButton) bool {
	if b, ok := gamepadButtonMapping[button]; ok {
		return w.currJoy[js].Button(b) != glfw.Press && w.prevJoy[js].Button(b) == glfw.Press
	}
	return false
}

// JoystickAxis returns the value of a joystick axis at the last call to Window.Update.
// If the axis index is out of range, this will return 0.
//
// This API is experimental.
func (w *Window) JoystickAxis(js pixel.Joystick, axis pixel.GamepadAxis) float64 {
	if a, ok := gamepadAxisMapping[axis]; ok {
		return float64(w.currJoy[js].Axis(a))
	}
	return 0
}

// Used internally during Window.UpdateInput to update the state of the joysticks.
func (w *Window) updateJoystickInput() {
	for js := pixel.Joystick1; js < pixel.Joystick(pixel.NumJoysticks); js++ {
		joystick, ok := joystickMapping[js]
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
