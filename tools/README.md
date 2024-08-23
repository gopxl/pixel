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
go run main.go bench run [names...]
```

Write benchmark stats to a file
```
go run main.go bench run [names...] -o my-stats.json
```

## Profiling
Run benchmark with cpu/mem profiling enabled
```
go run main.go bench run [names...] -c cpu.prof -m mem.prof
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

To add your own results to this file, create an entry in the [Machine Info](#machine-info) table with
a unique identifer and basic info about the computer where you are running the benchmarks.
On linux you can get most of the info from `lshw -short`. By default, benchmark stats will use the local username
from environment variables or the os package if `MACHINE_NAME` env is not provided.

Then run all benchmarks:
```
# Optional
export MACHINE_NAME=<machine-identifier>
export PIXEL_VERSION=<local-pixel-version>

go run main.go bench run --all
```

### Machine Info

| Machine            | OS/Distro           | CPU                           | Memory             | GPU            |
|--------------------|---------------------|-------------------------------|--------------------|----------------|
| bhperry-wsl        | Linux Ubuntu 20.04  | Intel i7-8086K @ 4.00GHz      | 8GiB               | RTX 2080       |
| bhperry-win10      | Windows 10          | Intel i7-8086K @ 4.00GHz      | 16GiB              | RTX 2080       |

### Stats

| Machine            | Pixel  | Benchmark                    | Duration | Frames | FPS Avg | FPS Min | FPS Max | FPS Stdev |
|--------------------|--------|------------------------------|----------|--------|---------|---------|---------|-----------|
| bhperry-wsl        | v2.2.1 | imdraw-moving                | 30s      | 2214   | 73.79   | 68      | 76      | 1.77      |
| bhperry-wsl        | v2.2.1 | imdraw-static                | 30s      | 2355   | 78.5    | 72      | 81      | 1.89      |
| bhperry-wsl        | v2.2.1 | sprite-moving                | 30.03s   | 1451   | 48.32   | 45      | 50      | 1.25      |
| bhperry-wsl        | v2.2.1 | sprite-moving-batched        | 30.01s   | 4085   | 136.12  | 127     | 142     | 3.17      |
| bhperry-wsl        | v2.2.1 | sprite-static                | 30.01s   | 1518   | 50.59   | 47      | 52      | 1.45      |
| bhperry-wsl        | v2.2.1 | sprite-static-batched        | 30.01s   | 5318   | 177.2   | 159     | 182     | 6.01      |
| bhperry-win10      | v2.2.1 | imdraw-moving                | 30.03s   | 1430   | 47.61   | 22      | 50      | 5.85      |
| bhperry-win10      | v2.2.1 | imdraw-static                | 30.02s   | 1569   | 52.27   | 51      | 53      | 0.64      |
| bhperry-win10      | v2.2.1 | sprite-moving                | 30.03s   | 1148   | 38.23   | 35      | 39      | 0.9       |
| bhperry-win10      | v2.2.1 | sprite-moving-batched        | 30s      | 39085  | 1302.79 | 1205    | 1329    | 23.93     |
| bhperry-win10      | v2.2.1 | sprite-static                | 30.04s   | 1218   | 40.54   | 38      | 42      | 0.88      |
| bhperry-win10      | v2.2.1 | sprite-static-batched        | 30s      | 40570  | 1352.29 | 1245    | 1380    | 26.04     |
