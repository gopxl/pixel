package pixel

type Button int

// String returns a human-readable string describing the Button.
func (b Button) String() string {
	name, ok := buttonNames[b]
	if !ok {
		return "Invalid"
	}
	return name
}

const ButtonUnknown Button = -1

const (
	// List of all mouse buttons.
	MouseButton1 Button = iota
	MouseButton2
	MouseButton3
	MouseButton4
	MouseButton5
	MouseButton6
	MouseButton7
	MouseButton8

	// Last iota
	// NOTE: These will be unexported in the future when Window is move to the pixel package.
	NumMouseButtons int = iota

	// Aliases
	MouseButtonLeft   = MouseButton1
	MouseButtonRight  = MouseButton2
	MouseButtonMiddle = MouseButton3
)

const (
	// List of all keyboard buttons.
	KeySpace = iota + Button(NumMouseButtons)
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDecimal
	KeyKPDivide
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPEnter
	KeyKPEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu

	// Last iota
	// NOTE: These will be unexported in the future when Window is move to the pixel package.
	NumKeys int = iota
)

var buttonNames = map[Button]string{
	ButtonUnknown:     "Unknown",
	MouseButton4:      "MouseButton4",
	MouseButton5:      "MouseButton5",
	MouseButton6:      "MouseButton6",
	MouseButton7:      "MouseButton7",
	MouseButton8:      "MouseButton8",
	MouseButtonLeft:   "MouseButtonLeft",
	MouseButtonRight:  "MouseButtonRight",
	MouseButtonMiddle: "MouseButtonMiddle",
	KeySpace:          "Space",
	KeyApostrophe:     "Apostrophe",
	KeyComma:          "Comma",
	KeyMinus:          "Minus",
	KeyPeriod:         "Period",
	KeySlash:          "Slash",
	Key0:              "0",
	Key1:              "1",
	Key2:              "2",
	Key3:              "3",
	Key4:              "4",
	Key5:              "5",
	Key6:              "6",
	Key7:              "7",
	Key8:              "8",
	Key9:              "9",
	KeySemicolon:      "Semicolon",
	KeyEqual:          "Equal",
	KeyA:              "A",
	KeyB:              "B",
	KeyC:              "C",
	KeyD:              "D",
	KeyE:              "E",
	KeyF:              "F",
	KeyG:              "G",
	KeyH:              "H",
	KeyI:              "I",
	KeyJ:              "J",
	KeyK:              "K",
	KeyL:              "L",
	KeyM:              "M",
	KeyN:              "N",
	KeyO:              "O",
	KeyP:              "P",
	KeyQ:              "Q",
	KeyR:              "R",
	KeyS:              "S",
	KeyT:              "T",
	KeyU:              "U",
	KeyV:              "V",
	KeyW:              "W",
	KeyX:              "X",
	KeyY:              "Y",
	KeyZ:              "Z",
	KeyLeftBracket:    "LeftBracket",
	KeyBackslash:      "Backslash",
	KeyRightBracket:   "RightBracket",
	KeyGraveAccent:    "GraveAccent",
	KeyWorld1:         "World1",
	KeyWorld2:         "World2",
	KeyEscape:         "Escape",
	KeyEnter:          "Enter",
	KeyTab:            "Tab",
	KeyBackspace:      "Backspace",
	KeyInsert:         "Insert",
	KeyDelete:         "Delete",
	KeyRight:          "Right",
	KeyLeft:           "Left",
	KeyDown:           "Down",
	KeyUp:             "Up",
	KeyPageUp:         "PageUp",
	KeyPageDown:       "PageDown",
	KeyHome:           "Home",
	KeyEnd:            "End",
	KeyCapsLock:       "CapsLock",
	KeyScrollLock:     "ScrollLock",
	KeyNumLock:        "NumLock",
	KeyPrintScreen:    "PrintScreen",
	KeyPause:          "Pause",
	KeyF1:             "F1",
	KeyF2:             "F2",
	KeyF3:             "F3",
	KeyF4:             "F4",
	KeyF5:             "F5",
	KeyF6:             "F6",
	KeyF7:             "F7",
	KeyF8:             "F8",
	KeyF9:             "F9",
	KeyF10:            "F10",
	KeyF11:            "F11",
	KeyF12:            "F12",
	KeyF13:            "F13",
	KeyF14:            "F14",
	KeyF15:            "F15",
	KeyF16:            "F16",
	KeyF17:            "F17",
	KeyF18:            "F18",
	KeyF19:            "F19",
	KeyF20:            "F20",
	KeyF21:            "F21",
	KeyF22:            "F22",
	KeyF23:            "F23",
	KeyF24:            "F24",
	KeyF25:            "F25",
	KeyKP0:            "KP0",
	KeyKP1:            "KP1",
	KeyKP2:            "KP2",
	KeyKP3:            "KP3",
	KeyKP4:            "KP4",
	KeyKP5:            "KP5",
	KeyKP6:            "KP6",
	KeyKP7:            "KP7",
	KeyKP8:            "KP8",
	KeyKP9:            "KP9",
	KeyKPDecimal:      "KPDecimal",
	KeyKPDivide:       "KPDivide",
	KeyKPMultiply:     "KPMultiply",
	KeyKPSubtract:     "KPSubtract",
	KeyKPAdd:          "KPAdd",
	KeyKPEnter:        "KPEnter",
	KeyKPEqual:        "KPEqual",
	KeyLeftShift:      "LeftShift",
	KeyLeftControl:    "LeftControl",
	KeyLeftAlt:        "LeftAlt",
	KeyLeftSuper:      "LeftSuper",
	KeyRightShift:     "RightShift",
	KeyRightControl:   "RightControl",
	KeyRightAlt:       "RightAlt",
	KeyRightSuper:     "RightSuper",
	KeyMenu:           "Menu",
}

// Joystick is a joystick or controller (gamepad).
type Joystick int

// String returns a human-readable string describing the Joystick.
func (j Joystick) String() string {
	name, ok := joystickNames[j]
	if !ok {
		return "Invalid"
	}
	return name
}

// List all of the joysticks.
const (
	Joystick1 Joystick = iota
	Joystick2
	Joystick3
	Joystick4
	Joystick5
	Joystick6
	Joystick7
	Joystick8
	Joystick9
	Joystick10
	Joystick11
	Joystick12
	Joystick13
	Joystick14
	Joystick15
	Joystick16

	// Last iota
	// NOTE: These will be unexported in the future when Window is move to the pixel package.
	NumJoysticks int = iota
)

var joystickNames = map[Joystick]string{
	Joystick1:  "Joystick1",
	Joystick2:  "Joystick2",
	Joystick3:  "Joystick3",
	Joystick4:  "Joystick4",
	Joystick5:  "Joystick5",
	Joystick6:  "Joystick6",
	Joystick7:  "Joystick7",
	Joystick8:  "Joystick8",
	Joystick9:  "Joystick9",
	Joystick10: "Joystick10",
	Joystick11: "Joystick11",
	Joystick12: "Joystick12",
	Joystick13: "Joystick13",
	Joystick14: "Joystick14",
	Joystick15: "Joystick15",
	Joystick16: "Joystick16",
}

// GamepadAxis corresponds to a gamepad axis.
type GamepadAxis int

// String returns a human-readable string describing the GamepadAxis.
func (ga GamepadAxis) String() string {
	name, ok := gamepadAxisNames[ga]
	if !ok {
		return "Invalid"
	}
	return name
}

// Gamepad axis IDs.
const (
	AxisLeftX GamepadAxis = iota
	AxisLeftY
	AxisRightX
	AxisRightY
	AxisLeftTrigger
	AxisRightTrigger

	// Last iota.
	// NOTE: These will be unexported in the future when Window is move to the pixel package.
	NumAxes int = iota
)

var gamepadAxisNames = map[GamepadAxis]string{
	AxisLeftX:        "AxisLeftX",
	AxisLeftY:        "AxisLeftY",
	AxisRightX:       "AxisRightX",
	AxisRightY:       "AxisRightY",
	AxisLeftTrigger:  "AxisLeftTrigger",
	AxisRightTrigger: "AxisRightTrigger",
}

// GamepadButton corresponds to a gamepad button.
type GamepadButton int

// String returns a human-readable string describing the GamepadButton.
func (gb GamepadButton) String() string {
	name, ok := gamepadButtonNames[gb]
	if !ok {
		return "Invalid"
	}
	return name
}

// Gamepad button IDs.
const (
	GamepadA GamepadButton = iota
	GamepadB
	GamepadX
	GamepadY
	GamepadLeftBumper
	GamepadRightBumper
	GamepadBack
	GamepadStart
	GamepadGuide
	GamepadLeftThumb
	GamepadRightThumb
	GamepadDpadUp
	GamepadDpadRight
	GamepadDpadDown
	GamepadDpadLeft

	// Last iota
	numGamepadButtons
	NumGamepadButtons = int(numGamepadButtons)

	// Aliases
	GamepadCross    = GamepadA
	GamepadCircle   = GamepadB
	GamepadSquare   = GamepadX
	GamepadTriangle = GamepadY
)

var gamepadButtonNames = map[GamepadButton]string{
	GamepadA:           "GamepadA",
	GamepadB:           "GamepadB",
	GamepadX:           "GamepadX",
	GamepadY:           "GamepadY",
	GamepadLeftBumper:  "GamepadLeftBumper",
	GamepadRightBumper: "GamepadRightBumper",
	GamepadBack:        "GamepadBack",
	GamepadStart:       "GamepadStart",
	GamepadGuide:       "GamepadGuide",
	GamepadLeftThumb:   "GamepadLeftThumb",
	GamepadRightThumb:  "GamepadRightThumb",
	GamepadDpadUp:      "GamepadDpadUp",
	GamepadDpadRight:   "GamepadDpadRight",
	GamepadDpadDown:    "GamepadDpadDown",
	GamepadDpadLeft:    "GamepadDpadLeft",
}
