//go:build js && wasm

package opengl

import (
	"strings"
	"syscall/js"

	"github.com/gopxl/pixel/v2"
)

// domCodeToButton maps KeyboardEvent.code values to pixel.Button constants.
// Only keys the game consumes are mapped; unknown codes yield ok=false.
var domCodeToButton = map[string]pixel.Button{
	"Space":        pixel.KeySpace,
	"Quote":        pixel.KeyApostrophe,
	"Comma":        pixel.KeyComma,
	"Minus":        pixel.KeyMinus,
	"Period":       pixel.KeyPeriod,
	"Slash":        pixel.KeySlash,
	"Digit0":       pixel.Key0,
	"Digit1":       pixel.Key1,
	"Digit2":       pixel.Key2,
	"Digit3":       pixel.Key3,
	"Digit4":       pixel.Key4,
	"Digit5":       pixel.Key5,
	"Digit6":       pixel.Key6,
	"Digit7":       pixel.Key7,
	"Digit8":       pixel.Key8,
	"Digit9":       pixel.Key9,
	"Semicolon":    pixel.KeySemicolon,
	"Equal":        pixel.KeyEqual,
	"KeyA":         pixel.KeyA,
	"KeyB":         pixel.KeyB,
	"KeyC":         pixel.KeyC,
	"KeyD":         pixel.KeyD,
	"KeyE":         pixel.KeyE,
	"KeyF":         pixel.KeyF,
	"KeyG":         pixel.KeyG,
	"KeyH":         pixel.KeyH,
	"KeyI":         pixel.KeyI,
	"KeyJ":         pixel.KeyJ,
	"KeyK":         pixel.KeyK,
	"KeyL":         pixel.KeyL,
	"KeyM":         pixel.KeyM,
	"KeyN":         pixel.KeyN,
	"KeyO":         pixel.KeyO,
	"KeyP":         pixel.KeyP,
	"KeyQ":         pixel.KeyQ,
	"KeyR":         pixel.KeyR,
	"KeyS":         pixel.KeyS,
	"KeyT":         pixel.KeyT,
	"KeyU":         pixel.KeyU,
	"KeyV":         pixel.KeyV,
	"KeyW":         pixel.KeyW,
	"KeyX":         pixel.KeyX,
	"KeyY":         pixel.KeyY,
	"KeyZ":         pixel.KeyZ,
	"BracketLeft":  pixel.KeyLeftBracket,
	"Backslash":    pixel.KeyBackslash,
	"BracketRight": pixel.KeyRightBracket,
	"Backquote":    pixel.KeyGraveAccent,
	"Escape":       pixel.KeyEscape,
	"Enter":        pixel.KeyEnter,
	"Tab":          pixel.KeyTab,
	"Backspace":    pixel.KeyBackspace,
	"Insert":       pixel.KeyInsert,
	"Delete":       pixel.KeyDelete,
	"ArrowRight":   pixel.KeyRight,
	"ArrowLeft":    pixel.KeyLeft,
	"ArrowDown":    pixel.KeyDown,
	"ArrowUp":      pixel.KeyUp,
	"PageUp":       pixel.KeyPageUp,
	"PageDown":     pixel.KeyPageDown,
	"Home":         pixel.KeyHome,
	"End":          pixel.KeyEnd,
	"CapsLock":     pixel.KeyCapsLock,
	"ScrollLock":   pixel.KeyScrollLock,
	"NumLock":      pixel.KeyNumLock,
	"PrintScreen":  pixel.KeyPrintScreen,
	"Pause":        pixel.KeyPause,
	"F1":           pixel.KeyF1,
	"F2":           pixel.KeyF2,
	"F3":           pixel.KeyF3,
	"F4":           pixel.KeyF4,
	"F5":           pixel.KeyF5,
	"F6":           pixel.KeyF6,
	"F7":           pixel.KeyF7,
	"F8":           pixel.KeyF8,
	"F9":           pixel.KeyF9,
	"F10":          pixel.KeyF10,
	"F11":          pixel.KeyF11,
	"F12":          pixel.KeyF12,
	"Numpad0":      pixel.KeyKP0,
	"Numpad1":      pixel.KeyKP1,
	"Numpad2":      pixel.KeyKP2,
	"Numpad3":      pixel.KeyKP3,
	"Numpad4":      pixel.KeyKP4,
	"Numpad5":      pixel.KeyKP5,
	"Numpad6":      pixel.KeyKP6,
	"Numpad7":      pixel.KeyKP7,
	"Numpad8":      pixel.KeyKP8,
	"Numpad9":      pixel.KeyKP9,
	"NumpadDecimal":  pixel.KeyKPDecimal,
	"NumpadDivide":   pixel.KeyKPDivide,
	"NumpadMultiply": pixel.KeyKPMultiply,
	"NumpadSubtract": pixel.KeyKPSubtract,
	"NumpadAdd":      pixel.KeyKPAdd,
	"NumpadEnter":    pixel.KeyKPEnter,
	"NumpadEqual":    pixel.KeyKPEqual,
	"ShiftLeft":      pixel.KeyLeftShift,
	"ControlLeft":    pixel.KeyLeftControl,
	"AltLeft":        pixel.KeyLeftAlt,
	"MetaLeft":       pixel.KeyLeftSuper,
	"ShiftRight":     pixel.KeyRightShift,
	"ControlRight":   pixel.KeyRightControl,
	"AltRight":       pixel.KeyRightAlt,
	"MetaRight":      pixel.KeyRightSuper,
	"ContextMenu":    pixel.KeyMenu,
}

// domMouseButton maps MouseEvent.button to the pixel.Button constants.
// Unmapped values yield ok=false.
var domMouseButton = map[int]pixel.Button{
	0: pixel.MouseButton1, // left
	1: pixel.MouseButton3, // middle
	2: pixel.MouseButton2, // right
	3: pixel.MouseButton4, // back
	4: pixel.MouseButton5, // forward
}

// initInput installs DOM keyboard and mouse listeners on the canvas. Events
// are translated into pixel.Button press/release/repeat events and mouse
// move/scroll events on the shared InputHandler; printable characters are
// appended to the Typed buffer.
func (w *Window) initInput() {
	keyDown := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		code := ev.Get("code").String()
		btn, ok := domCodeToButton[code]
		if !ok {
			return nil
		}

		if ev.Get("repeat").Bool() {
			w.input.ButtonEvent(btn, pixel.Repeat)
			w.fireButtonCallback(btn, pixel.Repeat)
		} else {
			w.input.ButtonEvent(btn, pixel.Press)
			w.fireButtonCallback(btn, pixel.Press)
		}

		// Feed printable characters into Typed buffer. Browsers give us the
		// localized string in `event.key`; only single-code-point values
		// correspond to real characters (others are "Enter", "Shift", etc).
		key := ev.Get("key").String()
		if r, ok := singleRune(key); ok {
			w.input.CharEvent(r)
			if w.charCallback != nil {
				w.charCallback(w, r)
			}
		}

		// Swallow default handling for game-consumed keys so the browser
		// doesn't scroll/tab-navigate while the canvas has focus.
		if shouldPreventDefault(code) {
			ev.Call("preventDefault")
		}
		return nil
	})

	keyUp := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		code := ev.Get("code").String()
		btn, ok := domCodeToButton[code]
		if !ok {
			return nil
		}
		w.input.ButtonEvent(btn, pixel.Release)
		w.fireButtonCallback(btn, pixel.Release)
		if shouldPreventDefault(code) {
			ev.Call("preventDefault")
		}
		return nil
	})

	// Blur clears all held keys so we don't get stuck-down buttons when the
	// tab loses focus mid-keypress.
	blur := js.FuncOf(func(this js.Value, args []js.Value) any {
		for _, btn := range domCodeToButton {
			w.input.ButtonEvent(btn, pixel.Release)
		}
		for _, btn := range domMouseButton {
			w.input.ButtonEvent(btn, pixel.Release)
		}
		// Clear touch state so no touches appear stuck after losing focus.
		for id := range w.activeTouches {
			delete(w.activeTouches, id)
		}
		if w.touchCount > 0 {
			w.touchCount = 0
			w.input.ButtonEvent(pixel.MouseButton1, pixel.Release)
		}
		return nil
	})

	mouseDown := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		btn, ok := domMouseButton[ev.Get("button").Int()]
		if !ok {
			return nil
		}
		// Prevent right-click context menu and middle-click autoscroll
		// while the canvas has focus.
		ev.Call("preventDefault")
		w.input.ButtonEvent(btn, pixel.Press)
		w.fireButtonCallback(btn, pixel.Press)
		// Ensure subsequent keydown events continue to land on the canvas.
		w.jsCanvas.Call("focus")
		return nil
	})

	mouseUp := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		btn, ok := domMouseButton[ev.Get("button").Int()]
		if !ok {
			return nil
		}
		w.input.ButtonEvent(btn, pixel.Release)
		w.fireButtonCallback(btn, pixel.Release)
		return nil
	})

	mouseMove := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		pos := w.mousePosFromEvent(args[0])
		w.input.MouseMoveEvent(pos)
		if w.mouseMovedCallback != nil {
			w.mouseMovedCallback(w, pos)
		}
		return nil
	})

	mouseEnter := js.FuncOf(func(this js.Value, args []js.Value) any {
		w.input.MouseEnteredEvent(true)
		if w.mouseEnteredCallback != nil {
			w.mouseEnteredCallback(w, true)
		}
		return nil
	})

	mouseLeave := js.FuncOf(func(this js.Value, args []js.Value) any {
		w.input.MouseEnteredEvent(false)
		if w.mouseEnteredCallback != nil {
			w.mouseEnteredCallback(w, false)
		}
		return nil
	})

	wheel := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		ev.Call("preventDefault")
		// Browsers deliver wheel deltas in pixels (deltaMode=0) or lines/
		// pages; we forward the raw pixel delta and let the caller scale it.
		// Y is inverted so scrolling up yields a positive value, matching
		// the desktop backend.
		dx := ev.Get("deltaX").Float()
		dy := -ev.Get("deltaY").Float()
		w.input.MouseScrollEvent(dx, dy)
		if w.scrollCallback != nil {
			w.scrollCallback(w, pixel.V(dx, dy))
		}
		return nil
	})

	contextMenu := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			args[0].Call("preventDefault")
		}
		return nil
	})

	// Touch events — map single-touch to MouseButton1 so the game receives
	// touch input on mobile browsers. preventDefault on all touch events stops
	// the browser's default pan/zoom/callout behaviours; passive:false is
	// required to allow that call.
	nonPassive := map[string]any{"passive": false}

	touchStart := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		ev.Call("preventDefault")
		changed := ev.Get("changedTouches")
		for i := 0; i < changed.Length(); i++ {
			t := changed.Index(i)
			id := t.Get("identifier").Int()
			pos := w.touchPosFromTouch(t)
			w.activeTouches[id] = pos
			w.input.MouseMoveEvent(pos)
		}
		// Fire Press only on the first finger down so MouseButton1 behaves as
		// a simple "any touch active" signal rather than toggling on each finger.
		if w.touchCount == 0 {
			w.input.ButtonEvent(pixel.MouseButton1, pixel.Press)
			w.fireButtonCallback(pixel.MouseButton1, pixel.Press)
		}
		w.touchCount += changed.Length()
		return nil
	})

	touchMove := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		ev.Call("preventDefault")
		// Update all moved touches in the active map.
		changed := ev.Get("changedTouches")
		for i := 0; i < changed.Length(); i++ {
			t := changed.Index(i)
			id := t.Get("identifier").Int()
			pos := w.touchPosFromTouch(t)
			if _, active := w.activeTouches[id]; active {
				w.activeTouches[id] = pos
			}
		}
		// Forward the first currently-active touch as the mouse position so
		// single-touch code (settings panel, etc.) continues to work.
		all := ev.Get("touches")
		if all.Length() > 0 {
			pos := w.touchPosFromTouch(all.Index(0))
			w.input.MouseMoveEvent(pos)
			if w.mouseMovedCallback != nil {
				w.mouseMovedCallback(w, pos)
			}
		}
		return nil
	})

	touchEnd := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		ev := args[0]
		ev.Call("preventDefault")
		changed := ev.Get("changedTouches")
		for i := 0; i < changed.Length(); i++ {
			t := changed.Index(i)
			id := t.Get("identifier").Int()
			if pos, ok := w.activeTouches[id]; ok {
				w.input.MouseMoveEvent(pos)
				delete(w.activeTouches, id)
			}
		}
		w.touchCount -= changed.Length()
		if w.touchCount <= 0 {
			w.touchCount = 0
			w.input.ButtonEvent(pixel.MouseButton1, pixel.Release)
			w.fireButtonCallback(pixel.MouseButton1, pixel.Release)
		}
		return nil
	})

	w.jsCanvas.Call("addEventListener", "keydown", keyDown)
	w.jsCanvas.Call("addEventListener", "keyup", keyUp)
	w.jsCanvas.Call("addEventListener", "blur", blur)
	w.jsCanvas.Call("addEventListener", "mousedown", mouseDown)
	w.jsCanvas.Call("addEventListener", "mouseup", mouseUp)
	w.jsCanvas.Call("addEventListener", "mousemove", mouseMove)
	w.jsCanvas.Call("addEventListener", "mouseenter", mouseEnter)
	w.jsCanvas.Call("addEventListener", "mouseleave", mouseLeave)
	w.jsCanvas.Call("addEventListener", "wheel", wheel, map[string]any{"passive": false})
	w.jsCanvas.Call("addEventListener", "contextmenu", contextMenu)
	w.jsCanvas.Call("addEventListener", "touchstart", touchStart, nonPassive)
	w.jsCanvas.Call("addEventListener", "touchmove", touchMove, nonPassive)
	w.jsCanvas.Call("addEventListener", "touchend", touchEnd, nonPassive)
	w.jsCanvas.Call("addEventListener", "touchcancel", touchEnd, nonPassive)
}

// mousePosFromEvent converts a MouseEvent's clientX/Y (in CSS pixels,
// relative to the viewport) into window-local coordinates in the backing
// store's pixel space. Y is flipped so the origin is at the bottom-left,
// matching the desktop backend.
func (w *Window) mousePosFromEvent(ev js.Value) pixel.Vec {
	rect := w.jsCanvas.Call("getBoundingClientRect")
	cssX := ev.Get("clientX").Float() - rect.Get("left").Float()
	cssY := ev.Get("clientY").Float() - rect.Get("top").Float()
	cssW := rect.Get("width").Float()
	cssH := rect.Get("height").Float()
	if cssW <= 0 || cssH <= 0 {
		return pixel.ZV
	}
	bounds := w.bounds
	x := bounds.Min.X + (cssX/cssW)*bounds.W()
	y := bounds.Min.Y + (1-cssY/cssH)*bounds.H()
	return pixel.V(x, y)
}

// touchPosFromTouch converts a single Touch object's clientX/Y into
// window-local pixel coordinates, using the same mapping as mousePosFromEvent.
func (w *Window) touchPosFromTouch(touch js.Value) pixel.Vec {
	rect := w.jsCanvas.Call("getBoundingClientRect")
	cssX := touch.Get("clientX").Float() - rect.Get("left").Float()
	cssY := touch.Get("clientY").Float() - rect.Get("top").Float()
	cssW := rect.Get("width").Float()
	cssH := rect.Get("height").Float()
	if cssW <= 0 || cssH <= 0 {
		return pixel.ZV
	}
	bounds := w.bounds
	x := bounds.Min.X + (cssX/cssW)*bounds.W()
	y := bounds.Min.Y + (1-cssY/cssH)*bounds.H()
	return pixel.V(x, y)
}

// fireButtonCallback invokes the user-registered button callback if any.
func (w *Window) fireButtonCallback(btn pixel.Button, action pixel.Action) {
	if w.buttonCallback != nil {
		w.buttonCallback(w, btn, action)
	}
}

func singleRune(s string) (rune, bool) {
	if s == "" {
		return 0, false
	}
	runes := []rune(s)
	if len(runes) != 1 {
		return 0, false
	}
	r := runes[0]
	if r < 0x20 || r == 0x7f {
		return 0, false
	}
	return r, true
}

// shouldPreventDefault reports whether the browser's default handling of a
// key should be suppressed while the canvas has focus.
func shouldPreventDefault(code string) bool {
	switch code {
	case "Space", "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight",
		"Tab", "Backspace", "Slash", "Quote":
		return true
	}
	if strings.HasPrefix(code, "F") && len(code) > 1 {
		// Let F12 (devtools) through; block other F-keys.
		if code == "F12" {
			return false
		}
		return true
	}
	return false
}
