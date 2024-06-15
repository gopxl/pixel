package internal

import "github.com/gopxl/pixel/v2"

type JoystickState struct {
	Connected [pixel.NumJoysticks]bool
	Name      [pixel.NumJoysticks]string
	Buttons   [pixel.NumJoysticks][]pixel.Action
	Axis      [pixel.NumJoysticks][]float32
}

// Returns if a button on a joystick is down, returning false if the button or joystick is invalid.
func (js *JoystickState) GetButton(joystick pixel.Joystick, button pixel.GamepadButton) bool {
	// Check that the joystick and button is valid, return false by default
	if js.Buttons[joystick] == nil || int(button) >= len(js.Buttons[joystick]) || button < 0 {
		return false
	}
	return js.Buttons[joystick][button] == pixel.Press
}

// Returns the value of a joystick axis, returning 0 if the button or joystick is invalid.
func (js *JoystickState) GetAxis(joystick pixel.Joystick, axis pixel.GamepadAxis) float64 {
	// Check that the joystick and axis is valid, return 0 by default.
	if js.Axis[joystick] == nil || int(axis) >= len(js.Axis[joystick]) || axis < 0 {
		return 0
	}
	return float64(js.Axis[joystick][axis])
}
