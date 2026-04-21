//go:build js && wasm

package opengl

import (
	"image/color"
	"syscall/js"

	"github.com/gopxl/glhf/v2"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/internal"
	"github.com/pkg/errors"
)

// CanvasElementID names the HTML canvas the WASM backend will attach to.
// Override before calling NewWindow to target a different element.
var CanvasElementID = "game"

// WindowConfig mirrors the desktop struct so call-sites compile unchanged.
// Most fields are no-ops under WASM; only Title and Bounds are honored.
type WindowConfig struct {
	Title                  string
	Icon                   []pixel.Picture
	Bounds                 pixel.Rect
	Position               pixel.Vec
	Monitor                *Monitor
	Resizable              bool
	Undecorated            bool
	NoIconify              bool
	AlwaysOnTop            bool
	TransparentFramebuffer bool
	VSync                  bool
	Maximized              bool
	Invisible              bool
	SamplesMSAA            int
	BoundsLimits           pixel.Rect
}

// MaxDevicePixelRatio caps the backing-store scale factor used by syncCanvasSize.
// iOS/WebKit drops the WebGL context when GPU memory is exhausted; DPR=3 on
// modern iPhones produces a ~12 MB framebuffer for a fullscreen canvas, which
// combined with game textures easily hits the iOS limit. DPR=2 gives sharp
// rendering at 2x scale for essentially no visual downgrade in a pixel-art game.
var MaxDevicePixelRatio = 2.0

// Window wraps an HTML5 canvas plus a WebGL2 rendering context. The internal
// Canvas handles all drawing; Update blits it to the default framebuffer and
// yields to requestAnimationFrame.
type Window struct {
	bounds pixel.Rect
	canvas *Canvas

	jsCanvas js.Value
	gl       js.Value

	closed        bool
	vsync         bool
	cursorVisible bool

	gpuRenderer  string // captured at init for diagnostics
	ctxLostCount int

	activeTouches map[int]pixel.Vec // keyed by Touch.identifier; updated by touch event handlers
	touchCount    int               // number of currently-held touches; guards MouseButton1 edge detection

	input                      internal.InputHandler
	prevJoy, currJoy, tempJoy  internal.JoystickState

	buttonCallback       func(win *Window, button pixel.Button, action pixel.Action)
	charCallback         func(win *Window, r rune)
	mouseEnteredCallback func(win *Window, entered bool)
	mouseMovedCallback   func(win *Window, pos pixel.Vec)
	scrollCallback       func(win *Window, scroll pixel.Vec)
}

var currWin *Window

// NewWindow binds to the configured HTML canvas, creates a WebGL2 context,
// and initializes glhf. Returns an error if the canvas or context is missing.
func NewWindow(cfg WindowConfig) (*Window, error) {
	doc := js.Global().Get("document")
	jsCanvas := doc.Call("getElementById", CanvasElementID)
	if !jsCanvas.Truthy() {
		return nil, errors.Errorf("canvas #%s not found", CanvasElementID)
	}

	gl := jsCanvas.Call("getContext", "webgl2", map[string]any{
		"alpha":              false,
		"antialias":          cfg.SamplesMSAA > 0,
		"premultipliedAlpha": true,
		"preserveDrawingBuffer": false,
	})
	if !gl.Truthy() {
		return nil, errors.New("webgl2 not available")
	}

	_, _, w, h := intBounds(cfg.Bounds)
	jsCanvas.Set("width", w)
	jsCanvas.Set("height", h)

	if cfg.Title != "" {
		doc.Set("title", cfg.Title)
	}

	glhf.SetContext(gl)
	glhf.Init()

	win := &Window{
		bounds:        cfg.Bounds,
		jsCanvas:      jsCanvas,
		gl:            gl,
		vsync:         cfg.VSync,
		cursorVisible: true,
		activeTouches: make(map[int]pixel.Vec),
	}

	// Capture GPU renderer string for diagnostics. WEBGL_debug_renderer_info
	// may be blocked by the browser (privacy), so fall back gracefully.
	if ext := gl.Call("getExtension", "WEBGL_debug_renderer_info"); ext.Truthy() {
		unmaskedRenderer := ext.Get("UNMASKED_RENDERER_WEBGL")
		if unmaskedRenderer.Truthy() {
			win.gpuRenderer = gl.Call("getParameter", unmaskedRenderer).String()
		}
	}
	if win.gpuRenderer == "" {
		win.gpuRenderer = gl.Call("getParameter", js.Global().Get("WebGL2RenderingContext").Get("RENDERER")).String()
	}

	win.canvas = NewCanvas(cfg.Bounds)
	currWin = win

	// Expose GPU renderer string to JS so the page can include it in crash
	// diagnostics (heartbeat, context-lost overlay, etc.).
	if win.gpuRenderer != "" {
		js.Global().Set("_gpuRenderer", win.gpuRenderer)
	}

	// Ensure the canvas can receive keyboard focus inside iframes.
	if jsCanvas.Get("tabIndex").Int() < 0 {
		jsCanvas.Set("tabIndex", 0)
	}
	jsCanvas.Call("focus")

	win.initInput()
	win.installContextLostHandler()

	return win, nil
}

// installContextLostHandler notifies the page when the WebGL context is lost.
// True in-place recovery is out of scope. If window.pixelOnContextLost is
// defined the page handles the response (show UI, prompt reload, etc.);
// otherwise we fall back to an immediate location.reload().
func (w *Window) installContextLostHandler() {
	lost := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			args[0].Call("preventDefault")
		}
		w.ctxLostCount++
		dpr := js.Global().Get("devicePixelRatio").Float()
		fbW := w.jsCanvas.Get("width").Int()
		fbH := w.jsCanvas.Get("height").Int()
		elapsed := js.Global().Get("performance").Call("now").Float()
		diag := js.Global().Get("Object").New()
		diag.Set("renderer", w.gpuRenderer)
		diag.Set("dpr", dpr)
		diag.Set("backingW", fbW)
		diag.Set("backingH", fbH)
		diag.Set("elapsedMs", elapsed)
		diag.Set("count", w.ctxLostCount)
		js.Global().Get("console").Call("warn", "webglcontextlost", diag)
		cb := js.Global().Get("pixelOnContextLost")
		if cb.Truthy() {
			cb.Invoke(diag)
		} else {
			js.Global().Get("location").Call("reload")
		}
		return nil
	})
	w.jsCanvas.Call("addEventListener", "webglcontextlost", lost)
}

// Destroy releases the WebGL context.
func (w *Window) Destroy() {
	w.closed = true
}

// Update blits the internal canvas to the default framebuffer and yields to
// the browser's requestAnimationFrame so the next frame can be scheduled.
func (w *Window) Update() {
	w.syncCanvasSize()
	w.SwapBuffers()
	w.UpdateInput()
	awaitAnimationFrame()
}

// syncCanvasSize updates the canvas backing store to match its current CSS size
// (times devicePixelRatio). Call before each frame so window resizes and
// fullscreen transitions propagate into window.Bounds() without a JS callback.
func (w *Window) syncCanvasSize() {
	cssW := w.jsCanvas.Get("clientWidth").Int()
	cssH := w.jsCanvas.Get("clientHeight").Int()
	if cssW <= 0 || cssH <= 0 {
		return
	}
	dpr := js.Global().Get("devicePixelRatio").Float()
	if dpr < 1 {
		dpr = 1
	} else if MaxDevicePixelRatio > 0 && dpr > MaxDevicePixelRatio {
		dpr = MaxDevicePixelRatio
	}
	targetW := int(float64(cssW) * dpr)
	targetH := int(float64(cssH) * dpr)
	curW := w.jsCanvas.Get("width").Int()
	curH := w.jsCanvas.Get("height").Int()
	if curW == targetW && curH == targetH {
		return
	}
	w.jsCanvas.Set("width", targetW)
	w.jsCanvas.Set("height", targetH)
	w.bounds = pixel.R(0, 0, float64(targetW), float64(targetH))
	w.canvas.SetBounds(w.bounds)
}

// SwapBuffers copies the current canvas texture into the default framebuffer.
func (w *Window) SwapBuffers() {
	fbW, fbH := w.framebufferSize()
	glhf.Bounds(0, 0, fbW, fbH)
	glhf.Clear(0, 0, 0, 0)

	w.canvas.gf.Frame().Begin()
	w.canvas.gf.Frame().Blit(
		nil,
		0, 0, w.canvas.Texture().Width(), w.canvas.Texture().Height(),
		0, 0, fbW, fbH,
	)
	w.canvas.gf.Frame().End()
}

func (w *Window) framebufferSize() (int, int) {
	return w.jsCanvas.Get("width").Int(), w.jsCanvas.Get("height").Int()
}

// awaitAnimationFrame blocks the calling goroutine until the browser fires
// the next rAF callback. This is how the WASM backend yields time to the
// JS event loop once per simulated frame.
func awaitAnimationFrame() {
	ch := make(chan struct{}, 1)
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) any {
		cb.Release()
		ch <- struct{}{}
		return nil
	})
	js.Global().Call("requestAnimationFrame", cb)
	<-ch
}

// Closed reports whether the window has been explicitly closed. The browser
// tab closing terminates the Go runtime, so this only reflects SetClosed.
func (w *Window) Closed() bool              { return w.closed }
func (w *Window) SetClosed(closed bool)     { w.closed = closed }

func (w *Window) Bounds() pixel.Rect        { return w.bounds }
func (w *Window) SetBounds(bounds pixel.Rect) {
	w.bounds = bounds
	_, _, width, height := intBounds(bounds)
	w.jsCanvas.Set("width", width)
	w.jsCanvas.Set("height", height)
	w.canvas.SetBounds(bounds)
}
func (w *Window) SetBoundsLimits(pixel.Rect) {}

func (w *Window) SetPos(pixel.Vec)    {}
func (w *Window) GetPos() pixel.Vec   { return pixel.ZV }
func (w *Window) SetTitle(title string) {
	js.Global().Get("document").Set("title", title)
}

func (w *Window) Focused() bool {
	return js.Global().Get("document").Call("hasFocus").Bool()
}

// ActiveTouches returns a snapshot of all currently-held touch positions in
// window-local pixel coordinates, one entry per active Touch.identifier.
// Use this instead of MousePosition for multi-touch hit-testing (e.g. virtual
// gamepads) because MousePosition only reflects the most-recent touch move.
func (w *Window) ActiveTouches() []pixel.Vec {
	out := make([]pixel.Vec, 0, len(w.activeTouches))
	for _, v := range w.activeTouches {
		out = append(out, v)
	}
	return out
}

func (w *Window) SetVSync(vsync bool) { w.vsync = vsync }
func (w *Window) VSync() bool         { return w.vsync }

func (w *Window) SetCursorVisible(visible bool) {
	w.cursorVisible = visible
	style := w.jsCanvas.Get("style")
	if visible {
		style.Set("cursor", "auto")
	} else {
		style.Set("cursor", "none")
	}
}
func (w *Window) SetCursorDisabled() {
	w.cursorVisible = false
	w.jsCanvas.Get("style").Set("cursor", "none")
}
func (w *Window) CursorVisible() bool { return w.cursorVisible }

func (w *Window) SetMonitor(*Monitor)      {}
func (w *Window) Monitor() *Monitor         { return nil }

// Note: must be called before any GL work. Kept for parity with desktop.
func (w *Window) begin() {
	if currWin != w {
		currWin = w
	}
}
func (w *Window) end() {}

// Drawing delegation -------------------------------------------------------

func (w *Window) MakeTriangles(t pixel.Triangles) pixel.TargetTriangles {
	return w.canvas.MakeTriangles(t)
}
func (w *Window) MakePicture(p pixel.Picture) pixel.TargetPicture {
	return w.canvas.MakePicture(p)
}
func (w *Window) SetMatrix(m pixel.Matrix)                 { w.canvas.SetMatrix(m) }
func (w *Window) SetColorMask(c color.Color)               { w.canvas.SetColorMask(c) }
func (w *Window) SetComposeMethod(cmp pixel.ComposeMethod) { w.canvas.SetComposeMethod(cmp) }
func (w *Window) SetSmooth(smooth bool)                    { w.canvas.SetSmooth(smooth) }
func (w *Window) Smooth() bool                             { return w.canvas.Smooth() }
func (w *Window) Clear(c color.Color)                      { w.canvas.Clear(c) }
func (w *Window) Color(at pixel.Vec) pixel.RGBA            { return w.canvas.Color(at) }
func (w *Window) Canvas() *Canvas                          { return w.canvas }

// Show/Hide/Focus are no-ops in a browser environment.
func (w *Window) Show()  {}
func (w *Window) Hide()  {}
func (w *Window) Focus() { w.jsCanvas.Call("focus") }

// Clipboard access goes through the browser's clipboard API when available.
// For now we return empty strings; writes are best-effort and ignored when
// the API is unavailable (e.g. insecure origin).
func (w *Window) ClipboardText() string        { return "" }
func (w *Window) SetClipboardText(text string) {}
func (w *Window) Clipboard() string            { return "" }
func (w *Window) SetClipboard(str string)      {}
