package debugger

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"os/exec"
)

var View *gocui.View

func Start() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrl2, gocui.ModNone, write); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrl3, gocui.ModNone, loop); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func loop(gui *gocui.Gui, view *gocui.View) error {
	v, _ := gui.View("shell")
	v.Write([]byte(`test`))
	return nil
}

func write(gui *gocui.Gui, view *gocui.View) error {
	v, _ := gui.View("shell")
	d := Debugger{
		View: v,
		Gui:  gui,
	}
	fmt.Fprintln(v, "hey")

	cmd := exec.Command("/bin/bash", "-c", "tmp/test.sh")
	cmd.Stdout = &d
	cmd.Stderr = &d

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	return nil
}

type Debugger struct {
	Gui  *gocui.Gui
	View *gocui.View
}

func (w Debugger) Write(p []byte) (n int, err error) {
	_, err = fmt.Fprintln(w.View, p)
	if err != nil {
		return 0, err
	}

	w.Gui.Update(func(gui *gocui.Gui) error {
		return nil
	})
	return len(p), nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	fourthX := maxX / 4
	halfX := maxX / 2

	if v, err := g.SetView("shell", 0, 0, maxX/2, maxY-20); err != nil {
		v.Title = "Shell"
		v.Wrap = true
		v.Editable = true
		View = v
	}

	if v, err := g.SetView("env", halfX+1, 0, halfX+(fourthX-1), maxY-20); err != nil {
		v.Title = "Environment"
		v.Editable = true
		fmt.Fprintln(v, "KEY: VALUE")
		fmt.Fprintln(v, "HELLO: TRUE")
		fmt.Fprintln(v, "PATH: /home")
		fmt.Fprintln(v, "USER: root")
	}

	if v, err := g.SetView("config", halfX+fourthX, 0, halfX+(fourthX+fourthX-1), maxY-20); err != nil {
		v.Title = "Configuration"
		v.Editable = true
		fmt.Fprintln(v, "Command: echo hello")
		fmt.Fprintln(v, "Working Directory: /home")
		fmt.Fprintln(v, "User: root")
		fmt.Fprintln(v, "Env: ")
		fmt.Fprintln(v, "USER: root")
	}
	return nil
}
