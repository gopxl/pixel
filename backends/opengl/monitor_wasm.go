//go:build js && wasm

package opengl

// Monitor represents a display. Under WASM we only ever have one logical
// monitor (the browser canvas host), so these are mostly stubs.
type Monitor struct{}

// VideoMode mirrors the desktop shape.
type VideoMode struct {
	Width       int
	Height      int
	RefreshRate int
}

func PrimaryMonitor() *Monitor    { return &Monitor{} }
func Monitors() []*Monitor        { return []*Monitor{{}} }
func (m *Monitor) Name() string   { return "Browser" }
func (m *Monitor) PhysicalSize() (width, height float64) {
	return 0, 0
}
func (m *Monitor) Position() (x, y float64) { return 0, 0 }
func (m *Monitor) Size() (width, height float64) {
	return 0, 0
}
func (m *Monitor) BitDepth() (red, green, blue int) {
	return 8, 8, 8
}
func (m *Monitor) RefreshRate() (rate float64) { return 60 }
func (m *Monitor) VideoModes() (vmodes []VideoMode) {
	return nil
}
