package pixel

func NewGamepad[Action, Axis, Button ~int, AxisValue comparable](name string) *Gamepad[Action, Axis, Button, AxisValue] {
	return &Gamepad[Action, Axis, Button, AxisValue]{
		name: name,
	}
}

type Gamepad[Action, Axis, Button ~int, AxisValue comparable] struct {
	connected bool
	name      string
	buttons   []Action
	axes      []AxisValue
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) Axis(axis Axis) AxisValue {
	return g.axes[axis]
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) Button(button Button) Action {
	if button >= Button(len(g.buttons)) || button < 0 {
		return -1
	}
	return g.buttons[button]
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) Connected() bool {
	return g.connected
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) Name() string {
	return g.name
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) NumAxes() int {
	return len(g.axes)
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) NumButtons() int {
	return len(g.buttons)
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) SetAxes(axes []AxisValue) {
	g.axes = axes
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) SetButtons(buttons []Action) {
	g.buttons = buttons
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) SetConnected(connected bool) {
	g.connected = connected
}

func (g *Gamepad[Action, Axis, Button, AxisValue]) SetName(name string) {
	g.name = name
}
