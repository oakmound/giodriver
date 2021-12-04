package main

import (
	"embed"
	"fmt"
	"image"
	"math/rand"
	"path/filepath"

	"github.com/oakmound/giodriver"
	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/scene"
)

const (
	minX = 0
	minY = 0
	maxX = 578
	maxY = 416
)

var cache = [360]*image.RGBA{}

func main() {
	oak.AddScene(
		"demo",
		scene.Scene{Start: func(*scene.Context) {
			render.Draw(render.NewDrawFPS(0.03, nil, 10, 10))
			render.Draw(render.NewLogicFPS(0.03, nil, 10, 20))

			instructionTxt := render.DefaultFont().NewText("Press (And Hold) K", 50, 20)
			render.Draw(instructionTxt, 0)

			layer := 0
			layerTxt := render.DefaultFont().NewIntText(&layer, 30, 20)
			layerTxt.SetLayer(100000000)
			render.Draw(layerTxt, 0)
			NewGopher(layer)
			layer++
			event.GlobalBind(event.Enter, func(event.CID, interface{}) int {
				if oak.IsDown("K") {
					NewGopher(layer)
					layer++
				}
				return 0
			})
			// Generate a rotation cache for comparison
			// Compare the use of the cache against the use of a reverting type below
			for i := 0; i < 360; i++ {
				s, err := render.LoadSprite(filepath.Join("assets", "images", "raw", "gopher11.png"))
				if err != nil {
					fmt.Println(err)
					return
				}
				s.Modify(mod.Rotate(float32(i)))
				cache[i] = s.GetRGBA()
			}
		},
		})

	render.SetDrawStack(
		render.NewCompositeR(),
	)

	oak.SetFS(assets)
	oak.Init("demo", func(c oak.Config) (oak.Config, error) {
		c.Driver = giodriver.Driver
		return c, nil
	})
}

//go:embed assets
var assets embed.FS

// Gopher is a basic bouncing renderable
type Gopher struct {
	*entities.Doodad
	deltaX, deltaY float64
	rotation       int
}

// Init sets up a gophers CID
func (g *Gopher) Init() event.CID {
	return event.NextID(g)
}

// NewGopher creates a gopher sprite to bounce around
func NewGopher(layer int) {
	goph := Gopher{}
	goph.Doodad = entities.NewDoodad(
		rand.Float64()*576,
		rand.Float64()*416,
		render.NewSwitch("goph", map[string]render.Modifiable{"goph": render.EmptyRenderable()}),
		//render.NewReverting(render.LoadSprite(filepath.Join("raw", "gopher11.png"))),
		goph.Init())
	goph.R.SetLayer(layer)
	goph.Bind("EnterFrame", gophEnter)
	goph.deltaX = 4 * float64(rand.Intn(2)*2-1)
	goph.deltaY = 4 * float64(rand.Intn(2)*2-1)
	goph.rotation = rand.Intn(360)
	render.Draw(goph.R, 0)
}

func gophEnter(cid event.CID, nothing interface{}) int {
	goph := event.GetEntity(cid).(*Gopher)

	// Compare against this version of rotation
	// (also swap the comments on lines in goph.Doodad's renderable)
	//goph.R.(*render.Reverting).RevertAndModify(1, render.Rotate(goph.rotation))
	goph.R.(*render.Switch).Add("goph", render.NewSprite(0, 0, cache[goph.rotation]))
	if goph.X() < minX || goph.X() > maxX {
		goph.deltaX *= -1
	}
	if goph.Y() < minY || goph.Y() > maxY {
		goph.deltaY *= -1
	}
	goph.SetPos(goph.deltaX+goph.X(), goph.deltaY+goph.Y())
	goph.rotation = (goph.rotation + 1) % 360
	return 0
}
