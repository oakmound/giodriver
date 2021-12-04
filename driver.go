package giodriver

import (
	"fmt"
	"image"
	"sync"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/shiny/screen"
)

var _ oak.Driver = Driver
var _ screen.Screen = &Screen{}
var _ screen.Image = &Image{}

var screenLock sync.Mutex
var screens []*Screen

type Screen struct {
	internalIndex       int
	lastWindowDestroyed chan (struct{})
}

func (s *Screen) NewImage(size image.Point) (screen.Image, error) {
	return Image{
		screen: s,
		size:   size,
		rgba:   image.NewRGBA(image.Rect(0, 0, size.X, size.Y)),
	}, nil
}

func (s *Screen) NewTexture(size image.Point) (screen.Texture, error) {
	return &Texture{
		screen: s,
		size:   size,
	}, nil
}

func (s *Screen) NewWindow(opts screen.WindowGenerator) (screen.Window, error) {
	if opts.Width == 0 || opts.Height == 0 {
		return nil, fmt.Errorf("invalid width/height: %d/%d", opts.Width, opts.Height)
	}
	window := &Window{
		screen: s,
		Deque:  &Deque{},
	}
	gioOpts := []app.Option{
		app.Size(unit.Value{
			V: float32(opts.Width),
			U: unit.UnitDp,
		}, unit.Value{
			V: float32(opts.Height),
			U: unit.UnitDp,
		}),
		app.Title(opts.Title),
	}
	if opts.Fullscreen {
		gioOpts = append(gioOpts, app.Fullscreen.Option())
	}
	window.gioWindow = app.NewWindow(gioOpts...)
	if opts.TopMost {
		window.gioWindow.Raise()
	}

	go window.handleEvents()

	return window, nil
}

func Driver(f func(screen.Screen)) {
	screenLock.Lock()
	screen := &Screen{
		lastWindowDestroyed: make(chan struct{}),
	}
	screens = append(screens, screen)
	screen.internalIndex = len(screens) - 1
	screenLock.Unlock()
	go f(screen)
	<-screen.lastWindowDestroyed
	screenLock.Lock()
	screens[screen.internalIndex] = nil
	screenLock.Unlock()
}
