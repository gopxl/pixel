package internal

import "github.com/gopxl/pixel/v2"

type InputState struct {
	Mouse   pixel.Vec
	Buttons [pixel.NumButtons]bool
	Repeat  [pixel.NumButtons]bool
	Scroll  pixel.Vec
	Typed   string
}

type InputHandler struct {
	Prev, Curr, temp InputState

	PressEvents, tempPressEvents     [pixel.NumButtons]bool
	ReleaseEvents, tempReleaseEvents [pixel.NumButtons]bool

	MouseInsideWindow bool
}

// SetMousePosition overrides the mouse position
// Called when the mouse is set to a point in the backend Window
func (ih *InputHandler) SetMousePosition(pos pixel.Vec) {
	ih.Prev.Mouse = pos
	ih.Curr.Mouse = pos
	ih.temp.Mouse = pos
}

// ButtonEvent sets the action state of a button for the next update
func (ih *InputHandler) ButtonEvent(button pixel.Button, action pixel.Action) {
	switch action {
	case pixel.Press:
		ih.tempPressEvents[button] = true
		ih.temp.Buttons[button] = true
	case pixel.Release:
		ih.tempReleaseEvents[button] = true
		ih.temp.Buttons[button] = false
	case pixel.Repeat:
		ih.temp.Repeat[button] = true
	}
}

// MouseMoveEvent sets the mouse position for the next update
func (ih *InputHandler) MouseMoveEvent(pos pixel.Vec) {
	ih.temp.Mouse = pos
}

// MouseScrollEvent adds to the scroll offset for the next update
func (ih *InputHandler) MouseScrollEvent(x, y float64) {
	ih.temp.Scroll.X += x
	ih.temp.Scroll.Y += y
}

// MouseEnteredEvent is called when the mouse enters or leaves the window
func (ih *InputHandler) MouseEnteredEvent(entered bool) {
	ih.MouseInsideWindow = entered
}

// CharEvent adds to the typed string for the next update
func (ih *InputHandler) CharEvent(r rune) {
	ih.temp.Typed += string(r)
}

func (ih *InputHandler) Update() {
	ih.Prev = ih.Curr
	ih.Curr = ih.temp

	ih.PressEvents = ih.tempPressEvents
	ih.ReleaseEvents = ih.tempReleaseEvents

	// Clear last frame's temporary status
	ih.tempPressEvents = [pixel.NumButtons]bool{}
	ih.tempReleaseEvents = [pixel.NumButtons]bool{}
	ih.temp.Repeat = [pixel.NumButtons]bool{}
	ih.temp.Scroll = pixel.ZV
	ih.temp.Typed = ""
}
