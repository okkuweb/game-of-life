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
	mouseAction mouseAction     // UI action
	interval  time.Duration // time interval between two frames
	pause bool
	options options
}

func main() {
	InitLogger()
	defer logFile.Close()
	opt := &options{width: 280, height: 65}
	gd := gruid.NewGrid(opt.width, opt.height)
	md := &model{grid: gd, pause: true, options: *opt}
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
	Type  actionType
}

type mouseAction struct {
	Type  actionType
	Location gruid.Point
}

type actionType int

const (
	ActionQuit   actionType = iota + 1
	ActionPause
	MouseMain
)

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	m.action = action{} // reset last action information
	switch msg := msg.(type) {
	case gruid.MsgInit:
		Log("Initializing")
		m.grid.Fill(gruid.Cell{Rune: ' '})
		return tick(m.interval)
	case timeMsg:
		Log("Pause: ", m.pause)
		if m.pause {
			break
		}
		return tick(m.interval + time.Millisecond * 100)
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	case gruid.MsgMouse:
		m.updateMouse(msg)
	}
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

	switch m.mouseAction.Type {
	case MouseMain:
		Log("Setting point")
		m.grid.Set(m.mouseAction.Location, gruid.Cell{Rune: '█'})
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

func (m *model) updateMouse(msg gruid.MsgMouse) {
	switch msg.Action {
	case gruid.MouseMain:
		m.mouseAction = mouseAction{Type: MouseMain, Location: msg.P}
	}
}

func (m *model) Draw() gruid.Grid {
	g2 := gruid.NewGrid(m.options.width, m.options.height)
	for p, c := range m.grid.All() {
		if !m.pause {
			m.AI(p, c, &g2)
		}
	}
	if !m.pause {
		m.grid = g2
	}
	return m.grid
}

func (m *model) AI(p gruid.Point, c gruid.Cell, g2 *gruid.Grid) gruid.Grid {
	around := gruid.NewRange(p.X-1, p.Y-1, p.X+2, p.Y+2)
	livecounter := 0
	for p2 := range around.Points() {
		if p2 == p || !m.grid.Contains(p2) {
			continue
		} else {
			c2 := m.grid.At(p2)
			if c2.Rune == '█' {
				livecounter++
			}
		}
	}
	if c.Rune == '█' { // If alive
		if livecounter == 2 || livecounter == 3 {
			g2.Set(p, gruid.Cell{Rune: '█'})
		} else {
			g2.Set(p, gruid.Cell{Rune: ' '})
		}
	} else { // If dead
		if livecounter == 3 {
			g2.Set(p, gruid.Cell{Rune: '█'})
		} else {
			g2.Set(p, gruid.Cell{Rune: ' '})
		}
	}
	return *g2
}

