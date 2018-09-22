package widget

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"strconv"
	"sync"
)

type Cursor struct {
	X, Y int
}

type Container struct {
	X, Y, Width, Height int
}

func cursor() func(format string, a ...interface{}) string {
	return color.New(color.FgBlack, color.BgWhite).SprintfFunc()
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Editor struct {
	Name       string
	Position   *Cursor
	Buffer     [][]string // TODO: find better structure
	Meta       Container
	WarpFactor int
	Mutex      *sync.Mutex
}

func (c *Cursor) Set(x, y int) {
	c.X = x
	c.Y = y
}

func (e *Editor) ResetWarpFactor() {
	e.WarpFactor = 1
}

func (e *Editor) Move(to Direction, by int) {
	e.Mutex.Lock()
	switch to {
	case Up:
		if e.Position.Y-by >= 0 {
			e.Position.Set(e.Position.X, e.Position.Y-by)
		}
	case Down:
		if e.Position.Y+by < e.Meta.Height {
			e.Position.Set(e.Position.X, e.Position.Y+by)
		}
	case Left:
		if e.Position.X-by >= 0 {
			e.Position.Set(e.Position.X-by, e.Position.Y)
		}
	case Right:
		if e.Position.X+by < e.Meta.Width {
			e.Position.Set(e.Position.X+by, e.Position.Y)
		}
	}
	e.Mutex.Unlock()
}

func createBuffer(with string, x, y int) [][]string {
	f := make([][]string, x)
	for a := range f {
		f[a] = make([]string, y)
		for b := range f[a] {
			f[a][b] = with
		}
	}
	return f
}

func NewEditor(name string, x, y, xm, ym int) *Editor {
	return &Editor{
		Name:       name,
		Buffer:     createBuffer("*", xm, ym),
		Position:   &Cursor{X: 0, Y: 0},
		Meta:       Container{X: x, Y: y, Width: xm + 2, Height: ym + 2},
		WarpFactor: 1,
		Mutex:      &sync.Mutex{},
	}
}

func (e *Editor) Layout(g *gocui.Gui) error {
	v, err := g.SetView(e.Name, e.Meta.X, e.Meta.Y, e.Meta.Width, e.Meta.Height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// we are good, lets write that buffer wont we
		//fmt.Fprintf(v, "%d%+v", e.WarpFactor, e.Position)
		for a := range e.Buffer {
			for ln := range e.Buffer[a] {
				line := strconv.Itoa(e.WarpFactor)
				if e.Position.Y == ln && e.Position.X == a {
					fmt.Fprintf(v, "%v", cursor()(line))
				} else {
					fmt.Fprintf(v, line)
				}
			}
			fmt.Fprintln(v)
		}
	}
	return nil
}
