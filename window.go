package giodriver

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/mouse"

	mobilekey "golang.org/x/mobile/event/key"
)

type Window struct {
	screen    *Screen
	gioWindow *app.Window
	toUpload  *image.RGBA
	*Deque
}

func (w *Window) Release()                                                                      {}
func (w *Window) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op)     {}
func (w *Window) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {}
func (w *Window) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op)       {}
func (w *Window) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	w.toUpload = src.(*Texture).rgba
	w.gioWindow.Invalidate()
}
func (w *Window) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {}
func (w *Window) Fill(dr image.Rectangle, src color.Color, op draw.Op)        {}

func (w *Window) Publish() screen.PublishResult {
	return screen.PublishResult{}
}

type SetAnimator interface {
	SetAnimating(bool)
}

func (w *Window) handleEvents() {
	var ops op.Ops
	for e := range w.gioWindow.Events() {
		switch e := e.(type) {
		case pointer.Event:
			//TODO: populate this event further
			w.Deque.Send(mouse.Event{
				X:         e.Position.X,
				Y:         e.Position.Y,
				Button:    mouse.Button(e.Buttons), // wrong
				Direction: 0,                       // ??
			})
		case key.Event:
			// TODO: handle more keys
			var r rune
			var code mobilekey.Code
			var dir mobilekey.Direction
			switch e.State {
			case key.Press:
				dir = mobilekey.DirPress
			case key.Release:
				dir = mobilekey.DirRelease
			default:
				fmt.Println("unknown dir", e.State)
			}
			switch e.Name {
			case "W":
				r = 'W'
				code = mobilekey.CodeW
			case "A":
				r = 'A'
				code = mobilekey.CodeA
			case "S":
				r = 'S'
				code = mobilekey.CodeS
			case "K":
				r = 'K'
				code = mobilekey.CodeK
			case "D":
				r = 'D'
				code = mobilekey.CodeD
			default:
				fmt.Println("unknown key", e.Name)
			}
			w.Deque.Send(mobilekey.Event{
				Rune:      r,
				Code:      code,
				Direction: dir,
			})
		case system.DestroyEvent:
			w.gioWindow.Close()
			// TODO: multi window
			w.screen.lastWindowDestroyed <- struct{}{}
			return
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			img := widget.Image{
				Src:   paint.NewImageOp(w.toUpload),
				Scale: 1,
			}
			img.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

// TODO: support screenopts