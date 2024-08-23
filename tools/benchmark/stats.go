package benchmark

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/user"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

var (
	machineName, pixelVersion string
)

func init() {
	machineName = getMachineName()
	pixelVersion = getPixelVersion()
}

// NewStats calculates statistics about a benchmark run
func NewStats(name string, duration time.Duration, frames int, frameSeconds []int) *Stats {
	stats := &Stats{
		Name:         name,
		Frames:       frames,
		Duration:     duration,
		Machine:      machineName,
		PixelVersion: pixelVersion,
	}

	milliseconds := stats.Duration.Milliseconds()
	if milliseconds > 0 {
		stats.AvgFPS = roundFloat(1000*float64(frames)/float64(milliseconds), 2)
	}

	fps := make([]float64, 0, len(frameSeconds))
	for i, frame := range frameSeconds {
		if i == 0 {
			fps = append(fps, float64(frame))
		} else {
			fps = append(fps, float64(frame-frameSeconds[i-1]))
		}
	}
	if len(fps) > 0 {
		stats.MinFPS = slices.Min(fps)
		stats.MaxFPS = slices.Max(fps)
		stats.StdevFPS = standardDeviation(fps)
	} else {
		// 1s or less test. Use average as a stand-in.
		stats.MinFPS = math.Floor(stats.AvgFPS)
		stats.MaxFPS = math.Ceil(stats.AvgFPS)
	}

	return stats
}

// Stats stores data about the performance of a benchmark run
type Stats struct {
	Name     string  `json:"name"`
	AvgFPS   float64 `json:"avgFPS"`
	MinFPS   float64 `json:"minFPS"`
	MaxFPS   float64 `json:"maxFPS"`
	StdevFPS float64 `json:"stdevFPS"`

	Frames   int           `json:"frames"`
	Duration time.Duration `json:"duration"`

	Machine      string `json:"machine"`
	PixelVersion string `json:"pixelVersion"`
}

// Print stats to stdout in a human-readable format
func (s *Stats) Print() {
	StatsCollection{s}.Print()
}

// StatsCollection holds stats from multiple benchmark runs
type StatsCollection []*Stats

func (sc StatsCollection) Print() {
	data := make([][]string, len(sc))
	for i, stats := range sc {
		data[i] = []string{
			stats.Machine,
			stats.PixelVersion,
			stats.Name,
			roundDuration(stats.Duration, 2).String(),
			toString(stats.Frames),
			toString(stats.AvgFPS),
			toString(stats.MinFPS),
			toString(stats.MaxFPS),
			toString(stats.StdevFPS),
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Machine", "Pixel", "Benchmark", "Duration", "Frames", "FPS Avg", "FPS Min", "FPS Max", "FPS Stdev"}
	widths := map[string]int{
		"Machine":   18,
		"Pixel":     6,
		"Benchmark": 28,
	}
	for i, header := range headers {
		minWidth := widths[header]
		if minWidth == 0 {
			minWidth = 6
		}
		table.SetColMinWidth(i, minWidth)
	}
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
}

// Dump writes a JSON file of all stored statistics to the given path
func (sc StatsCollection) Dump(path string) error {
	bytes, err := json.Marshal(sc)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, bytes, 0666); err != nil {
		return err
	}
	return nil
}

// roundFloat rounds the value to the given number of decimal places
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// roundDuration rounds the duration to the given number of decimal places based on the unit
func roundDuration(duration time.Duration, precision uint) time.Duration {
	durationRounding := time.Duration(math.Pow(10, float64(precision)))
	switch {
	case duration > time.Second:
		return duration.Round(time.Second / durationRounding)
	case duration > time.Millisecond:
		return duration.Round(time.Millisecond / durationRounding)
	case duration > time.Microsecond:
		return duration.Round(time.Microsecond / durationRounding)
	default:
		return duration
	}
}

func toString(val any) string {
	switch v := val.(type) {
	case float64:
		return fmt.Sprintf("%v", roundFloat(v, 2))
	case float32:
		return fmt.Sprintf("%v", roundFloat(float64(v), 2))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// standardDeviation calulates the variation of the given values relative to the average
func standardDeviation(values []float64) float64 {
	var sum, avg, stdev float64
	for _, val := range values {
		sum += val
	}
	count := float64(len(values))
	avg = sum / count

	for _, val := range values {
		stdev += math.Pow(val-avg, 2)
	}
	stdev = math.Sqrt(stdev / count)
	return stdev
}

func getMachineName() string {
	envs := []string{"MACHINE_NAME", "USER", "USERNAME"}
	var name string
	for _, env := range envs {
		name = os.Getenv(env)
		if name != "" {
			return name
		}
	}
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return "unknown"
}

func getPixelVersion() string {
	ver := os.Getenv("PIXEL_VERSION")
	if ver != "" {
		return ver
	}

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if dep.Path == "github.com/gopxl/pixel/v2" {
				return strings.Split(dep.Version, "-")[0]
			}
		}
	}
	return "x.y.z"
}
