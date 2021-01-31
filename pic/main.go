package main

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func init() {
	os.Chdir("/home/mokee/repos/gosound/pic")
}

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Blah",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var pts [1024]float64
	im := imdraw.New(nil)
	im.Color = colornames.Red

	m := NewMonitor()
	for !win.Closed() {
		m.NewFrame()

		for i := range pts {
			x := 10 * 2 * math.Pi * (float64(i)/float64(len(pts)) - 0.5*float64(m.T)/float64(time.Second))
			pts[i] = 20 * math.Sin(x)
		}

		win.Clear(colornames.White)
		im.Clear()
		for i, pt := range pts {
			im.Push(pixel.V(float64(i), win.Bounds().H()/2+pt))
			im.Circle(2, 0)
		}
		im.Draw(win)

		win.SetTitle("Hello | " + strconv.Itoa(int(m.FPS())) + " FPS")
		win.Update()
	}
}

type Monitor struct {
	lastFrame time.Time
	T         time.Duration
	DT        time.Duration

	Frame        int
	timePerFrame float64
}

func NewMonitor() Monitor {
	return Monitor{lastFrame: time.Now()}
}

func (m *Monitor) NewFrame() {
	m.Frame++

	now := time.Now()
	m.DT = now.Sub(m.lastFrame)
	m.lastFrame = now
	m.T += m.DT

	if m.Frame > 1 {
		m.timePerFrame += 0.001 * (m.DT.Seconds() - m.timePerFrame)
	}
}

func (m *Monitor) FPS() float64 {
	return 1 / m.timePerFrame
}
