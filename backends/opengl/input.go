package opengl

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gopxl/mainthread/v2"
	"github.com/gopxl/pixel/v2"
)

// Pressed returns whether the Button is currently pressed down.
func (w *Window) Pressed(button pixel.Button) bool {
	return w.input.Curr.Buttons[button]
}

// JustPressed returns whether the Button has been pressed in the last frame.
func (w *Window) JustPressed(button pixel.Button) bool {
	return w.input.PressEvents[button]
}

// JustReleased returns whether the Button has been released in the last frame.
func (w *Window) JustReleased(button pixel.Button) bool {
	return w.input.ReleaseEvents[button]
}

// Repeated returns whether a repeat event has been triggered on button.
//
// Repeat event occurs repeatedly when a button is held down for some time.
func (w *Window) Repeated(button pixel.Button) bool {
	return w.input.Curr.Repeat[button]
}

// MousePosition returns the current mouse position in the Window's Bounds.
func (w *Window) MousePosition() pixel.Vec {
	return w.input.Curr.Mouse
}

// MousePreviousPosition returns the previous mouse position in the Window's Bounds.
func (w *Window) MousePreviousPosition() pixel.Vec {
	return w.input.Prev.Mouse
}

// SetMousePosition positions the mouse cursor anywhere within the Window's Bounds.
func (w *Window) SetMousePosition(v pixel.Vec) {
	mainthread.Call(func() {
		if (v.X >= 0 && v.X <= w.bounds.W()) &&
			(v.Y >= 0 && v.Y <= w.bounds.H()) {
			w.window.SetCursorPos(
				v.X+w.bounds.Min.X,
				(w.bounds.H()-v.Y)+w.bounds.Min.Y,
			)
			w.input.SetMousePosition(v)
		}
	})
}

// MouseInsideWindow returns true if the mouse position is within the Window's Bounds.
func (w *Window) MouseInsideWindow() bool {
	return w.input.MouseInsideWindow
}

// MouseScroll returns the mouse scroll amount (in both axes) since the last call to Window.Update.
func (w *Window) MouseScroll() pixel.Vec {
	return w.input.Curr.Scroll
}

func (w *Window) MousePreviousScroll() pixel.Vec {
	return w.input.Prev.Scroll
}

// Typed returns the text typed on the keyboard since the last call to Window.Update.
func (w *Window) Typed() string {
	return w.input.Curr.Typed
}

var actionMapping = map[glfw.Action]pixel.Action{
	glfw.Release: pixel.Release,
	glfw.Press:   pixel.Press,
	glfw.Repeat:  pixel.Repeat,
}

var mouseButtonMapping = map[glfw.MouseButton]pixel.Button{
	glfw.MouseButton1: pixel.MouseButton1,
	glfw.MouseButton2: pixel.MouseButton2,
	glfw.MouseButton3: pixel.MouseButton3,
	glfw.MouseButton4: pixel.MouseButton4,
	glfw.MouseButton5: pixel.MouseButton5,
	glfw.MouseButton6: pixel.MouseButton6,
	glfw.MouseButton7: pixel.MouseButton7,
	glfw.MouseButton8: pixel.MouseButton8,
}

var keyButtonMapping = map[glfw.Key]pixel.Button{
	glfw.KeyUnknown:      pixel.UnknownButton,
	glfw.KeySpace:        pixel.KeySpace,
	glfw.KeyApostrophe:   pixel.KeyApostrophe,
	glfw.KeyComma:        pixel.KeyComma,
	glfw.KeyMinus:        pixel.KeyMinus,
	glfw.KeyPeriod:       pixel.KeyPeriod,
	glfw.KeySlash:        pixel.KeySlash,
	glfw.Key0:            pixel.Key0,
	glfw.Key1:            pixel.Key1,
	glfw.Key2:            pixel.Key2,
	glfw.Key3:            pixel.Key3,
	glfw.Key4:            pixel.Key4,
	glfw.Key5:            pixel.Key5,
	glfw.Key6:            pixel.Key6,
	glfw.Key7:            pixel.Key7,
	glfw.Key8:            pixel.Key8,
	glfw.Key9:            pixel.Key9,
	glfw.KeySemicolon:    pixel.KeySemicolon,
	glfw.KeyEqual:        pixel.KeyEqual,
	glfw.KeyA:            pixel.KeyA,
	glfw.KeyB:            pixel.KeyB,
	glfw.KeyC:            pixel.KeyC,
	glfw.KeyD:            pixel.KeyD,
	glfw.KeyE:            pixel.KeyE,
	glfw.KeyF:            pixel.KeyF,
	glfw.KeyG:            pixel.KeyG,
	glfw.KeyH:            pixel.KeyH,
	glfw.KeyI:            pixel.KeyI,
	glfw.KeyJ:            pixel.KeyJ,
	glfw.KeyK:            pixel.KeyK,
	glfw.KeyL:            pixel.KeyL,
	glfw.KeyM:            pixel.KeyM,
	glfw.KeyN:            pixel.KeyN,
	glfw.KeyO:            pixel.KeyO,
	glfw.KeyP:            pixel.KeyP,
	glfw.KeyQ:            pixel.KeyQ,
	glfw.KeyR:            pixel.KeyR,
	glfw.KeyS:            pixel.KeyS,
	glfw.KeyT:            pixel.KeyT,
	glfw.KeyU:            pixel.KeyU,
	glfw.KeyV:            pixel.KeyV,
	glfw.KeyW:            pixel.KeyW,
	glfw.KeyX:            pixel.KeyX,
	glfw.KeyY:            pixel.KeyY,
	glfw.KeyZ:            pixel.KeyZ,
	glfw.KeyLeftBracket:  pixel.KeyLeftBracket,
	glfw.KeyBackslash:    pixel.KeyBackslash,
	glfw.KeyRightBracket: pixel.KeyRightBracket,
	glfw.KeyGraveAccent:  pixel.KeyGraveAccent,
	glfw.KeyWorld1:       pixel.KeyWorld1,
	glfw.KeyWorld2:       pixel.KeyWorld2,
	glfw.KeyEscape:       pixel.KeyEscape,
	glfw.KeyEnter:        pixel.KeyEnter,
	glfw.KeyTab:          pixel.KeyTab,
	glfw.KeyBackspace:    pixel.KeyBackspace,
	glfw.KeyInsert:       pixel.KeyInsert,
	glfw.KeyDelete:       pixel.KeyDelete,
	glfw.KeyRight:        pixel.KeyRight,
	glfw.KeyLeft:         pixel.KeyLeft,
	glfw.KeyDown:         pixel.KeyDown,
	glfw.KeyUp:           pixel.KeyUp,
	glfw.KeyPageUp:       pixel.KeyPageUp,
	glfw.KeyPageDown:     pixel.KeyPageDown,
	glfw.KeyHome:         pixel.KeyHome,
	glfw.KeyEnd:          pixel.KeyEnd,
	glfw.KeyCapsLock:     pixel.KeyCapsLock,
	glfw.KeyScrollLock:   pixel.KeyScrollLock,
	glfw.KeyNumLock:      pixel.KeyNumLock,
	glfw.KeyPrintScreen:  pixel.KeyPrintScreen,
	glfw.KeyPause:        pixel.KeyPause,
	glfw.KeyF1:           pixel.KeyF1,
	glfw.KeyF2:           pixel.KeyF2,
	glfw.KeyF3:           pixel.KeyF3,
	glfw.KeyF4:           pixel.KeyF4,
	glfw.KeyF5:           pixel.KeyF5,
	glfw.KeyF6:           pixel.KeyF6,
	glfw.KeyF7:           pixel.KeyF7,
	glfw.KeyF8:           pixel.KeyF8,
	glfw.KeyF9:           pixel.KeyF9,
	glfw.KeyF10:          pixel.KeyF10,
	glfw.KeyF11:          pixel.KeyF11,
	glfw.KeyF12:          pixel.KeyF12,
	glfw.KeyF13:          pixel.KeyF13,
	glfw.KeyF14:          pixel.KeyF14,
	glfw.KeyF15:          pixel.KeyF15,
	glfw.KeyF16:          pixel.KeyF16,
	glfw.KeyF17:          pixel.KeyF17,
	glfw.KeyF18:          pixel.KeyF18,
	glfw.KeyF19:          pixel.KeyF19,
	glfw.KeyF20:          pixel.KeyF20,
	glfw.KeyF21:          pixel.KeyF21,
	glfw.KeyF22:          pixel.KeyF22,
	glfw.KeyF23:          pixel.KeyF23,
	glfw.KeyF24:          pixel.KeyF24,
	glfw.KeyF25:          pixel.KeyF25,
	glfw.KeyKP0:          pixel.KeyKP0,
	glfw.KeyKP1:          pixel.KeyKP1,
	glfw.KeyKP2:          pixel.KeyKP2,
	glfw.KeyKP3:          pixel.KeyKP3,
	glfw.KeyKP4:          pixel.KeyKP4,
	glfw.KeyKP5:          pixel.KeyKP5,
	glfw.KeyKP6:          pixel.KeyKP6,
	glfw.KeyKP7:          pixel.KeyKP7,
	glfw.KeyKP8:          pixel.KeyKP8,
	glfw.KeyKP9:          pixel.KeyKP9,
	glfw.KeyKPDecimal:    pixel.KeyKPDecimal,
	glfw.KeyKPDivide:     pixel.KeyKPDivide,
	glfw.KeyKPMultiply:   pixel.KeyKPMultiply,
	glfw.KeyKPSubtract:   pixel.KeyKPSubtract,
	glfw.KeyKPAdd:        pixel.KeyKPAdd,
	glfw.KeyKPEnter:      pixel.KeyKPEnter,
	glfw.KeyKPEqual:      pixel.KeyKPEqual,
	glfw.KeyLeftShift:    pixel.KeyLeftShift,
	glfw.KeyLeftControl:  pixel.KeyLeftControl,
	glfw.KeyLeftAlt:      pixel.KeyLeftAlt,
	glfw.KeyLeftSuper:    pixel.KeyLeftSuper,
	glfw.KeyRightShift:   pixel.KeyRightShift,
	glfw.KeyRightControl: pixel.KeyRightControl,
	glfw.KeyRightAlt:     pixel.KeyRightAlt,
	glfw.KeyRightSuper:   pixel.KeyRightSuper,
	glfw.KeyMenu:         pixel.KeyMenu,
}

func (w *Window) SetButtonCallback(callback func(win *Window, button pixel.Button, action pixel.Action)) {
	w.buttonCallback = callback
}

func (w *Window) SetCharCallback(callback func(win *Window, r rune)) {
	w.charCallback = callback
}

func (w *Window) SetMouseEnteredCallback(callback func(win *Window, entered bool)) {
	w.mouseEnteredCallback = callback
}

func (w *Window) SetMouseMovedCallback(callback func(win *Window, pos pixel.Vec)) {
	w.mouseMovedCallback = callback
}

func (w *Window) SetScrollCallback(callback func(win *Window, scroll pixel.Vec)) {
	w.scrollCallback = callback
}

func (w *Window) initInput() {
	mainthread.Call(func() {
		w.window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
			if b, buttonOk := mouseButtonMapping[button]; buttonOk {
				if a, actionOk := actionMapping[action]; actionOk {
					w.input.ButtonEvent(b, a)
					if w.buttonCallback != nil {
						w.buttonCallback(w, b, a)
					}
				}
			}
		})

		w.window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeyUnknown {
				return
			}
			if b, buttonOk := keyButtonMapping[key]; buttonOk {
				if a, actionOk := actionMapping[action]; actionOk {
					w.input.ButtonEvent(b, a)
					if w.buttonCallback != nil {
						w.buttonCallback(w, b, a)
					}
				}
			}
		})

		w.window.SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
			if entered && w.cursor != nil {
				w.window.SetCursor(w.cursor)
			}
			w.input.MouseEnteredEvent(entered)
			if w.mouseEnteredCallback != nil {
				w.mouseEnteredCallback(w, entered)
			}
		})

		w.window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
			pos := pixel.V(
				x+w.bounds.Min.X,
				(w.bounds.H()-y)+w.bounds.Min.Y,
			)
			w.input.MouseMoveEvent(pos)
			if w.mouseMovedCallback != nil {
				w.mouseMovedCallback(w, pos)
			}
		})

		w.window.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
			w.input.MouseScrollEvent(xoff, yoff)
			if w.scrollCallback != nil {
				w.scrollCallback(w, pixel.V(xoff, yoff))
			}
		})

		w.window.SetCharCallback(func(_ *glfw.Window, r rune) {
			w.input.CharEvent(r)
			if w.charCallback != nil {
				w.charCallback(w, r)
			}
		})
	})
}

// UpdateInput polls window events. Call this function to poll window events
// without swapping buffers. Note that the Update method invokes UpdateInput.
func (w *Window) UpdateInput() {
	mainthread.Call(func() {
		glfw.PollEvents()
	})
	w.doUpdateInput()
}

// UpdateInputWait blocks until an event is received or a timeout. If timeout is 0
// then it will wait indefinitely
func (w *Window) UpdateInputWait(timeout time.Duration) {
	mainthread.Call(func() {
		if timeout <= 0 {
			glfw.WaitEvents()
		} else {
			glfw.WaitEventsTimeout(timeout.Seconds())
		}
	})
	w.doUpdateInput()
}

// internal input bookkeeping
func (w *Window) doUpdateInput() {
	w.input.Update()
	w.updateJoystickInput()
}
