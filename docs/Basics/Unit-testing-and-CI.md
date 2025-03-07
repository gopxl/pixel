[Previous Tutorial](./Using-a-custom-fragment-shader.md)

This tutorial shows how to write unit tests that work
despite the all windows must run in the OpenGL main thread,
and how to run them in a CI.

## TestMain

To be able to write unit tests that use an OpenGL window,
all test must run in the same OpenGL context.
This can be achieved by putting a file `main_test.go` like this one into each package that contains tests:

```go
package foo_test

import (
	"os"
	"testing"

	"github.com/gopxl/pixel/v2/backends/opengl"
)

func TestMain(m *testing.M) {
	opengl.Run(func() {
		os.Exit(m.Run())
	})
}
```

## Continuous integration (CI)

A CI (like Github Actions) does usually not provide a display.
Tests that use an OpenGL window must therefore be run using [`xvfb`](https://en.wikipedia.org/wiki/Xvfb).

Here is a minimal Github Actions example:

```yml
name: Tests

on:
  push:
    branches:
    - main
  pull_request:

jobs:
  test:
    name: Run tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.24.x'
      - name: Install dependencies
        run: |
          sudo apt-get update -y
          sudo apt-get install -y libgl1-mesa-dev xorg-dev xvfb
          go get ./...
      - name: Run tests
        run: |
          xvfb-run go test -v ./...
```