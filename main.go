package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
)

type options struct {
	width int
	height int
}

type model struct {
	grid   gruid.Grid
	action action
	interval  time.Duration
	pause bool
	options options
	ui      *ui.Label
	frame   gruid.Grid
}

func main() {
	InitLogger()
	defer logFile.Close()
	opt := &options{width: 280, height: 65}
	gd := gruid.NewGrid(opt.width, opt.height)
	md := &model{grid: gd, pause: true, options: *opt}

	st := gruid.Style{}
	md.ui = &ui.Label{
		Box:     &ui.Box{Title: ui.Text(" Game of Life ")},
		Content: ui.StyledText{}.WithStyle(st.WithFg(ColorGreen)),
	}

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
	Location gruid.Point
}

type actionType int

const (
	MouseMain   actionType = iota + 1
	ActionQuit
	ActionPause
)

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	m.action = action{} // reset last action information
	switch msg := msg.(type) {
	case gruid.MsgInit:
		Log("Initializing")
		m.frame = gruid.NewGrid(m.options.width, m.options.height)
		m.grid.Fill(gruid.Cell{Rune: ' '})
		m.frame.Fill(gruid.Cell{Rune: ' '})
		return tick(m.interval)
	case timeMsg:
		m.ui.SetText("Pause: " + strconv.FormatBool(m.pause))
		if m.pause {
			break
		}
		g2 := gruid.NewGrid(m.options.width, m.options.height)
		if !m.pause {
			for p, c := range m.frame.All() {
				m.AI(p, c, &g2)
			}
		}
		m.frame = g2
		return tick(m.interval + time.Millisecond * 200)
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

func (m *model) AI(p gruid.Point, c gruid.Cell, g2 *gruid.Grid) gruid.Grid {
	around := gruid.NewRange(p.X-1, p.Y-1, p.X+2, p.Y+2)
	livecounter := 0
	for p2 := range around.Points() {
		if p2 == p || !m.frame.Contains(p2) {
			continue
		} else {
			c2 := m.frame.At(p2)
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

func (m *model) handleAction() gruid.Effect {

	switch m.action.Type {
	case ActionPause:
		m.pause = !m.pause
		if !m.pause {
			return tick(m.interval)
		}
	case ActionQuit:
		return gruid.End()
	case MouseMain:
		m.frame.Set(m.action.Location, gruid.Cell{Rune: '█'})
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
		m.action = action{Type: MouseMain, Location: msg.P}
	}
}

func (m *model) Draw() gruid.Grid {
	m.grid.Copy(m.frame)
	m.ui.Draw(m.grid.Slice(gruid.NewRange(0, 0, 20, 5)))
	return m.grid
}

