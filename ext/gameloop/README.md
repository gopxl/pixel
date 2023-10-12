# Plugin - Gameloop

A simple plugin that allows you to turn any `pixel window` into a managable entity that can be used in a game loop.

Defines an `EasyWindow` interface with the following methods

```go
type EasyWindow interface {
	Win() *opengl.Window // get underlying GLFW window
	Setup() error         // setup window
	Update() error        // update window
	Draw() error          // draw to window
}
```

The loop in pseudo code looks like this:

```go

win.Setup()

for !win.Closed() {
    win.Update()
    win.Draw()
}
```

Define your `Update` logic to handle user input and update the state of your game. Define your `Draw` logic to draw the state of your game to the window.

The game loop has the ability to handle multiple windows as well. 

## Example

```go
package main

window1 := MyNewWindow() // assume MyNewWindow implements EasyWindow interface
window2 := MyOtherWindow() // assume MyOtherWindow implements EasyWindow interface

manager := NewWindowManager()
manager.InsertWindows([]opengl.EasyWindow{
    window1,
    window2,
})

manager.SetFPS(60) // set the FPS of the game loop

if err := manager.Loop(); err != nil {
    panic(err)
}
```

The above code will allow you to run a game loop with multiple windows. The manager *assumes* that the first window given to the `InsertWindows` method is the main window, and will 
block main loop until that main window is closed. The manager will then close all other windows and exit the game loop.

