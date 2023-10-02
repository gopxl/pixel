package pixelgl

import "time"

var frameTick *time.Ticker

func setFPS(fps int) {
	if fps <= 0 {
		frameTick = nil
	} else {
		ms := 1.0 / float64(fps) * 1000.0
		dur := time.Duration(ms) * time.Millisecond
		frameTick = time.NewTicker(dur)
	}
}

type EasyWindow interface {
	Win() *Window  // get underlying GLFW window
	Update() error // update window
	Draw() error   // draw to window
}

type WindowManager struct {
	Windows []EasyWindow
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

func (wm *WindowManager) Loop(fps float64) error {
	// assumes first index is main loop
	win := wm.Windows[0].Win()

	for !win.Closed() {
		if err := wm.update(); err != nil {
			return err
		}
		if err := wm.draw(); err != nil {
			return err
		}
	}

	// update GFLW window
	for _, win := range wm.Windows {
		win.Win().Update()
	}

	if frameTick != nil {
		<-frameTick.C
	}

	return nil
}
