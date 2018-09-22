// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/jroimartin/gocui"
	"godray/widget"
	"log"
)

func setupKeys(e *widget.Editor, g *gocui.Gui) error {
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			e.Position.Set(5, 5)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, Encapsulate(MoveDown, e)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			e.Move(widget.Left, e.WarpFactor)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			e.Move(widget.Right, e.WarpFactor)
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	//xm, ym := g.Size()
	editor := widget.NewEditor("main", 1, 1, 10, 10)
	g.SetManager(editor)
	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		editor.Move(widget.Down, 1)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func Encapsulate(f func(*gocui.Gui, *gocui.View, *widget.Editor) error, with *widget.Editor) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f(g, v, with)
	}
}

func MoveDown(g *gocui.Gui, v *gocui.View, e *widget.Editor) error {
	e.Move(widget.Down, e.WarpFactor)
	return nil
}
