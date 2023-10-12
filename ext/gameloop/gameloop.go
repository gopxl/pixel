package gameloop

import (
	"errors"
	"time"

	"github.com/gopxl/pixel/v2/backends/opengl"
)

type EasyWindow interface {
	Win() *pixelgl.Window // get underlying GLFW window
	Setup() error         // setup window
	Update() error        // update window
	Draw() error          // draw to window
}

type WindowManager struct {
	Windows        []EasyWindow
	currentFps     float64
	targetDuration time.Duration
}

func NewWindowManager() *WindowManager {
	return &WindowManager{}
}

func (wm *WindowManager) SetFPS(fps int) error {
	if fps <= 0 {
		return errors.New("FPS must be greater than 0")
	}
	us := 1.0 / float64(fps) * 1000000.0
	wm.targetDuration = time.Duration(us) * time.Microsecond
	return nil
}

func (wm *WindowManager) FPS() float64 {
	return wm.currentFps
}
func (wm *WindowManager) InsertWindow(win EasyWindow) error {
	wm.Windows = append(wm.Windows, win)
	return nil
}

func (wm *WindowManager) InsertWindows(wins []EasyWindow) error {
	for _, win := range wins {
		if err := wm.InsertWindow(win); err != nil {
			return err
		}
	}
	return nil
}

func (wm *WindowManager) update() error {
	for _, win := range wm.Windows {
		if err := win.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (wm *WindowManager) draw() error {
	for _, win := range wm.Windows {
		if err := win.Draw(); err != nil {
			return err
		}
	}
	return nil
}

func (wm *WindowManager) Loop() error {
	// assumes first index is main loop
	win := wm.Windows[0].Win()

	if win == nil {
		panic("no main window")
	}

	// setup windows
	for _, win := range wm.Windows {
		if err := win.Setup(); err != nil {
			return err
		}
	}

	for !win.Closed() {
		start := time.Now()

		if err := wm.update(); err != nil {
			return err
		}
		if err := wm.draw(); err != nil {
			return err
		}

		// update GFLW window
		for _, win := range wm.Windows {
			win.Win().Update()
		}

		// calculate FPS
		elapsed := time.Since(start)

		if elapsed < wm.targetDuration {
			time.Sleep(wm.targetDuration - elapsed)
		}
		wm.currentFps = 1000000.0 / float64(time.Since(start).Microseconds())
	}
	return nil
}
