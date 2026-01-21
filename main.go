package main

import (
	"context"
	"fmt"
	"time"

	"codeberg.org/anaseto/gruid"
)

type options struct {
	width int
	height int
}

type model struct {
	grid   gruid.Grid // drawing grid
	action action     // UI action
	interval  time.Duration // time interval between two frames
	pause bool
}

func main() {
	InitLogger()
	defer logFile.Close()
	// TODO: Move this grid stuff to a grid file
	opt := &options{width: 80, height: 24}
	gd := gruid.NewGrid(opt.width, opt.height)
	md := &model{grid: gd}
	initDriver()

	app := gruid.NewApp(gruid.AppConfig{
		Driver: driver,
		Model:  md,
	})

	if err := app.Start(context.Background()); err != nil {
		fmt.Println(err)
	}
}

type action struct {
	Type  actionType  // kind of action (movement, quitting, ...)
	Delta gruid.Point // direction for ActionMovement
}

type actionType int

const (
	ActionQuit   actionType = iota + 1
	ActionPause
)

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	m.action = action{} // reset last action information
	switch msg := msg.(type) {
	case gruid.MsgInit:
		Log("Initializing")
		return tick(m.interval)
	case timeMsg:
		Log("Pause: ", m.pause)
		if m.pause {
			break
		}
		return tick(m.interval + time.Millisecond * 500)
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	}
	// Handle action (if any).
	return m.handleAction()
}

type timeMsg time.Time

func tick(d time.Duration) gruid.Cmd {
	t := time.NewTimer(d)
	return func() gruid.Msg {
		return timeMsg(<-t.C)
	}
}

func (m *model) handleAction() gruid.Effect {
	switch m.action.Type {
	case ActionPause:
		m.pause = !m.pause
		if !m.pause {
			return tick(m.interval)
		}
	case ActionQuit:
		return gruid.End()
	}
	return nil
}

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {
	switch msg.Key {
	case gruid.KeySpace, "p", "P":
		m.action = action{Type: ActionPause}
	case gruid.KeyEscape, "Q":
		m.action = action{Type: ActionQuit}
	}
}

var drawswitch bool
func (m *model) Draw() gruid.Grid {
	if drawswitch == false {
		Log("false")
		m.grid.Fill(gruid.Cell{Rune: '#'})
		drawswitch = true
	} else {
		Log("true")
		m.grid.Fill(gruid.Cell{Rune: ' '})
		drawswitch = false
	}
	return m.grid
}
