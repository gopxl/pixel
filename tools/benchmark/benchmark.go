package benchmark

import (
	"fmt"
	"slices"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

var Benchmarks = &Registry{benchmarks: map[string]Config{}}

// Config defines how to run a given benchmark, along with metadata describing it
type Config struct {
	Name        string
	Description string

	// New returns the benchmark to be executed
	New func(win *opengl.Window) (Benchmark, error)
	// Duration sets the maximum duration to run the benchmark
	Duration time.Duration
	// WindowConfig defines the input parameters to the benchmark's window
	WindowConfig opengl.WindowConfig
}

// Run executes the benchmark and calculates statistics about its performance
func (c Config) Run() (*Stats, error) {
	fmt.Printf("Running benchmark %s\n", c.Name)

	windowConfig := c.WindowConfig
	title := windowConfig.Title
	if title == "" {
		title = c.Name
	}
	windowConfig.Title = fmt.Sprintf("%s | FPS -", title)

	if windowConfig.Bounds.Empty() {
		windowConfig.Bounds = pixel.R(0, 0, 1024, 1024)
	}
	if windowConfig.Position.Eq(pixel.ZV) {
		windowConfig.Position = pixel.V(50, 50)
	}

	duration := c.Duration
	if duration == 0 {
		duration = 10 * time.Second
	}

	win, err := opengl.NewWindow(windowConfig)
	if err != nil {
		return nil, err
	}
	defer win.Destroy()

	benchmark, err := c.New(win)
	if err != nil {
		return nil, err
	}

	frame := 0
	frameSeconds := make([]int, 0)
	prevFrameCount := 0
	second := time.NewTicker(time.Second)
	done := time.NewTicker(duration)
	start := time.Now()
	last := start
loop:
	for frame = 0; !win.Closed(); frame++ {
		now := time.Now()
		benchmark.Step(win, now.Sub(last).Seconds())
		last = now
		win.Update()

		select {
		case <-second.C:
			frameSeconds = append(frameSeconds, frame)
			win.SetTitle(fmt.Sprintf("%s | FPS %v", title, frame-prevFrameCount))
			prevFrameCount = frame
		case <-done.C:
			break loop
		default:
		}
	}
	stats := NewStats(c.Name, time.Since(start), frame, frameSeconds)

	if win.Closed() {
		return nil, fmt.Errorf("window closed early")
	}

	return stats, err
}

// Benchmark provides hooks into the stages of a window's lifecycle
type Benchmark interface {
	Step(win *opengl.Window, delta float64)
}

// Registry is a collection of benchmark configs
type Registry struct {
	benchmarks map[string]Config
}

// List returns a copy of all registered benchmark configs
func (r *Registry) List() []Config {
	configs := make([]Config, len(r.benchmarks))
	for i, name := range r.ListNames() {
		configs[i] = r.benchmarks[name]
		i++
	}
	return configs
}

// ListNames returns a sorted list of all registered benchmark names
func (r *Registry) ListNames() []string {
	names := make([]string, len(r.benchmarks))
	i := 0
	for name := range r.benchmarks {
		names[i] = name
		i++
	}
	slices.Sort(names)
	return names
}

// Add a benchmark config to the registry
func (r *Registry) Add(configs ...Config) {
	for _, config := range configs {
		r.benchmarks[config.Name] = config
	}
}

// Get a benchmark config by name
func (r *Registry) Get(name string) (Config, error) {
	config, ok := r.benchmarks[name]
	if !ok {
		return config, fmt.Errorf("unknown benchmark %s", name)
	}

	return config, nil
}
