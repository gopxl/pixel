# Benchmarking

The `bench` command provides a set of tools used for benchmarking the performance of pixel under various scenarios.
It is intended to be a development tool for comparing the performance of new implementations in pixel against previous iterations.

## Usage

List available benchmarks
```
go run main.go bench ls
```

Run a benchmark
```
go run main.go bench run imdraw-static
```

Write benchmark stats to a file
```
go run main.go bench run imdraw-static -o imdraw-static-stats.json
```

## Profiling
Run benchmark with cpu/mem profiling enabled
```
go run main.go bench run -c cpu.prof -m mem.prof
```

View profile on cmdline
```
go tool pprof cpu.prof
```

View profile in browser (requires [graphviz](https://graphviz.org/download/))
```
go tool pprof -http :9000 cpu.prof
```

## Results

### Machine Info

Information about the machines used to record benchmark stats

| Machine            | OS/Distro           | CPU                           | Memory             | GPU            |
|--------------------|---------------------|-------------------------------|--------------------|----------------|
| bhperry-wsl        | Linux Ubuntu 20.04  | Intel i7-8086K @ 4.00GHz      | 8GiB               | RTX 2080       |
| bhperry-win10      | Windows 10          | Intel i7-8086K @ 4.00GHz      | 16GiB              | RTX 2080       |

### Stats

| Machine            | Pixel  | Benchmark                    | Duration | Frames | FPS Avg | FPS Min | FPS Max | FPS Stdev |
|--------------------|--------|------------------------------|----------|--------|---------|---------|---------|-----------|
| bhperry-wsl        | v2.2.0 | imdraw-moving                | 30.01s   | 2232   | 74.37   | 60      | 78      | 3.45      |
| bhperry-wsl        | v2.2.0 | imdraw-static                | 30.02s   | 2334   | 77.75   | 73      | 80      | 1.2       |
| bhperry-wsl        | v2.2.0 | sprite-moving                | 30.03s   | 1452   | 48.35   | 45      | 50      | 1.05      |
| bhperry-wsl        | v2.2.0 | sprite-moving-batched        | 30.01s   | 4004   | 133.42  | 127     | 139     | 2.45      |
| bhperry-wsl        | v2.2.0 | sprite-static                | 30.02s   | 1534   | 51.1    | 48      | 52      | 0.91      |
| bhperry-wsl        | v2.2.0 | sprite-static-batched        | 30s      | 5293   | 176.43  | 163     | 179     | 2.99      |
| bhperry-win10      | v2.2.0 | imdraw-moving                | 30.03s   | 1425   | 47.45   | 21      | 49      | 4.96      |
| bhperry-win10      | v2.2.0 | imdraw-static                | 30s      | 1533   | 51.1    | 50      | 52      | 0.55      |
| bhperry-win10      | v2.2.0 | sprite-moving                | 30.02s   | 1145   | 38.15   | 37      | 39      | 0.46      |
| bhperry-win10      | v2.2.0 | sprite-moving-batched        | 30s      | 39753  | 1325.06 | 1269    | 1348    | 15.1      |
| bhperry-win10      | v2.2.0 | sprite-static                | 30.01s   | 1214   | 40.45   | 40      | 41      | 0.5       |
| bhperry-win10      | v2.2.0 | sprite-static-batched        | 30s      | 39513  | 1317.06 | 1299    | 1336    | 10.1      |
