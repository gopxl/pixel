Normally, it is sufficient to grab the Go MSI installer from the website in
order to set up the toolchain. However, some packages that provide Go wrappers
for C libraries rely on [cgo](https://golang.org/cmd/cgo/) tool, which in turn,
needs the GCC toolchain in order to build the glue code. Also, 3rd-party
dependencies are usually hosted on services like GitHub, thus
[Git](https://git-scm.com/) is also needed. This mini-guide illustrates how to
setup a convenient development environment on Windows using MSYS2.

# Go
Follow the official instructions on Go
[website](https://golang.org/doc/install#windows) and install it, if you haven't
done this yet. The default installation directory is usually fine for most
users. Afterwards, locate the `go.exe` executable in the `bin` folder of the
installation root and annotate the path to it. Typical path usually looks like
`C:/Go/bin`.

# MSYS2
The [MSYS2](http://www.msys2.org/) project aims to provide an isolated
POSIX-like environment and a versatile package manager for easy software
configuration and building on Windows.

Please follow the guide on the website and after having completed the setup,
open the _MSYS2 MSYS_ shell from Windows menu, which will spawn a terminal with
a Bash shell.

MSYS2 uses [Pacman](https://wiki.archlinux.org/index.php/pacman) as package
manager for automated software installation and we will use it to bootstrap our
development environment.

# GCC toolchain
The GCC compiler suite and the development libraries needed for _cgo_ can be
installed with just one command:

    pacman -S --needed base-devel mingw-w64-i686-toolchain mingw-w64-x86_64-toolchain

Then add it to `$PATH`:

    echo 'export PATH=/mingw64/bin:$PATH' >> ~/.bashrc

# Git
There are two options available for Git installation: either from the
[website](https://git-scm.com) or via _pacman_. If you have Git already
installed via official installer, locate the path to `git.exe` executable in the
installation path and annotate it. Otherwise, install it via _pacman_:

    pacman -S git

# Environment setup
The `PATH` and `GOPATH` environment variables need to be propery configured to
point to Go and Git executables. We'll append their initialization to user's
local `.bashrc` file, which holds Bash shell customization.

_NOTE_: All paths *should* be specified in UNIX filesystem notation, although
this is not a strict requirement, since Bash translates Windows paths to Unix
automatically, but it helps keep consistency and prevents some nasty situations.
In MSYS2 environment, drive letters (C:, D:, etc) start from root (`/`), for
example: `C:\Windows` translates to `/c/Windows`. Backslash (`\`) is used for
spaces escaping. Case does not matter, alhtough it's preferable to maintain it.

## Go tools path
Recall the path to `go.exe` binary and append it to `.bashrc`, if the default Go
installation path was used, the command will look like this:

    echo 'export PATH=/c/Go/bin:$PATH' >> ~/.bashrc

## Git path (optional)
This is needed only if Git was installed via official installer, as MSYS2 does
not know where to locate the `git.exe` binary. For a standard Git installation
the command will be:

    echo 'export PATH=/c/Program\ Files/Git/bin:$PATH' >> ~/.bashrc

## GOPATH (optional)
This variable specifies the root for Go package sources and binaries. Since Go
1.8 it is not mandatory anymore, and on Windows it defaults to
`%USERPROFILE%\Go`, which is normally a path like `C:\Users\username\Go`. If you
need it to point to a different location, append the override command to
`.bashrc`, for example:

     echo 'export GOPATH=$USERPROFILE/go' >> ~/.bashrc

# Building Go packages
From now on, all development will be done in the just configured MinGW shell.
Note that this is not the MSYS2 shell we just used! The former is where
development occurs, like building C/C++ or Go programs for Windows, the latter
is used for the subsystem management, package installation, etc. Although both
shells share the same `.bashrc`, the MinGW shell has additional configuration,
and for example, picks the right toolchain based on the chosen architecture.
From Windows menu, open the MinGW shell from the launcher provided by MSYS2, for
the x64 architecture it is called _MSYS2 MinGW 64-bit_.

Now we can build and execute Go packages in an environment which closely
resembles a typical UNIX one.

The following example checks out and builds the platformer demo of the
[Pixel](https://github.com/gopxl/pixel) 2D game development framework. Its
hardware-accelerated backend uses [go-gl](https://github.com/go-gl/gl) package,
which provides a Go wrapper for OpenGL C libraries.

    go get github.com/gopxl/pixel-examples/v2
    cd $(go env GOPATH)/src/github.com/gopxl/pixel-examples/platformer
    go get
    go build
    ./platformer.exe
