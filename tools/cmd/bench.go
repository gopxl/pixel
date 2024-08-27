package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"text/tabwriter"
	"time"

	"github.com/gopxl/pixel/tools/benchmark"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/spf13/cobra"
)

var (
	benchRunAll bool
	benchRunOutput,
	benchRunCpuprofile,
	benchRunMemprofile string
	benchRunDuration time.Duration

	benchStatsInput string
)

func NewBenchCmd() *cobra.Command {
	bench := &cobra.Command{
		Use:   "bench",
		Short: "Benchmark the pixel library",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	bench.AddCommand(NewBenchLsCmd(), NewBenchRunCmd(), NewBenchStatsCmd())
	return bench
}

func NewBenchRunCmd() *cobra.Command {
	run := &cobra.Command{
		Use:   "run [names...] [opts]",
		Short: "Run one or more benchmark tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			if benchRunAll {
				args = benchmark.Benchmarks.ListNames()
			} else if len(args) == 0 {
				return fmt.Errorf("requires at least one benchmark")
			}
			cmd.SilenceUsage = true

			// Start CPU profile
			if benchRunCpuprofile != "" {
				f, err := os.Create(benchRunCpuprofile)
				if err != nil {
					return fmt.Errorf("could not create CPU profile: %v", err)
				}
				defer f.Close()
				if err := pprof.StartCPUProfile(f); err != nil {
					return fmt.Errorf("could not start CPU profile: %v", err)
				}
				defer pprof.StopCPUProfile()
			}

			// Run benchmark(s)
			benchStats := make(benchmark.StatsCollection, len(args))
			var err error
			run := func() {
				var config benchmark.Config
				for i, name := range args {
					config, err = benchmark.Benchmarks.Get(name)
					if err != nil {
						return
					}

					if benchRunDuration != 0 {
						config.Duration = benchRunDuration
					}

					var stats *benchmark.Stats
					stats, err = config.Run()
					if err != nil {
						return
					}
					benchStats[i] = stats
				}
			}

			opengl.Run(run)
			if err != nil {
				return err
			}
			fmt.Println()
			benchStats.Print()

			// Dump memory profile
			if benchRunMemprofile != "" {
				f, err := os.Create(benchRunMemprofile)
				if err != nil {
					return fmt.Errorf("could not create memory profile: %v", err)
				}
				defer f.Close()
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(f); err != nil {
					return fmt.Errorf("could not write memory profile: %v", err)
				}
			}

			// Dump stats
			if benchRunOutput != "" {
				err := benchStats.Dump(benchRunOutput)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	run.Flags().BoolVarP(&benchRunAll, "all", "a", false, "Run all registered benchmarks")
	run.Flags().StringVarP(&benchRunOutput, "output", "o", "", "Output path for statistics file")
	run.Flags().DurationVarP(&benchRunDuration, "duration", "d", 0, "Override duration for benchmark runs")
	run.Flags().StringVarP(&benchRunCpuprofile, "cpuprofile", "c", "", "CPU profiling file")
	run.Flags().StringVarP(&benchRunMemprofile, "memprofile", "m", "", "Memory profiling file")
	return run
}

func NewBenchLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List available benchmarks",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 1, 4, 8, ' ', 0)
			for _, config := range benchmark.Benchmarks.List() {
				fmt.Fprintf(w, "%s\t%s\n", config.Name, config.Description)
			}
			w.Flush()
		},
	}
}

func NewBenchStatsCmd() *cobra.Command {
	stats := &cobra.Command{
		Use:          "stats -i [path/to/stats.json]",
		Short:        "Pretty print the contents of a stats file",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := os.ReadFile(benchStatsInput)
			if err != nil {
				return err
			}

			var benchStats benchmark.StatsCollection
			if err := json.Unmarshal(bytes, &benchStats); err != nil {
				return err
			}
			benchStats.Print()

			return nil
		},
	}

	stats.Flags().StringVarP(&benchStatsInput, "input", "i", "", "Input path for statistics file")
	stats.MarkFlagRequired("input")
	return stats
}
