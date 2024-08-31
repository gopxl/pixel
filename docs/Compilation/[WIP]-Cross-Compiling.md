Besides setting a `GOOS` and/or `GOARCH` at compile time, you'll need to specify a `CC` in order to compile the c-based libraries.

MingW is recommended when the target platform is windows: https://www.mingw-w64.org/downloads/

An example build command for your project might look like the following: 

`CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build ./cmd/main/main.go -o bin/project.exe`
