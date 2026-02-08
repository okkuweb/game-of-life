package main

import (
	"context"
	"fmt"
	"time"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
)

type options struct {
	width  int
	height int
	speed  int
}

type model struct {
	grid       gruid.Grid
	action     action
	heldAction actionType
	interval   time.Duration
	pause      bool
	opts       options
	ui         *ui.Label
	frame      gruid.Grid
	entities   map[gruid.Point]bool
}

func main() {
	InitLogger()
	defer logFile.Close()
	Log("Starting game")
	opts := &options{width: 80, height: 24, speed: 200}
	gd := gruid.NewGrid(opts.width, opts.height)
	entities := make(map[gruid.Point]bool)
	md := &model{grid: gd, pause: true, opts: *opts, entities: entities}

	md.ui = &ui.Label{
		Box:     &ui.Box{Title: ui.Text(" Game of Life ")},
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
	Type     actionType
	Location gruid.Point
	Update   updateType
}

type updateType int

const (
	Map updateType = iota + 1
	UI
)

type actionType int

const (
	MouseMain actionType = iota + 1
	MouseSecondary
	MouseRelease
	MouseMove
	ActionQuit
	ActionPause
	ActionSpeedUp
	ActionSpeedDown
	ActionEnlargeMapY
	ActionShrinkMapY
	ActionEnlargeMapX
	ActionShrinkMapX
)

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	m.action = action{} // reset last action information
	switch msg := msg.(type) {
	case gruid.MsgInit:
		m.frame = gruid.NewGrid(m.opts.width, m.opts.height)
		m.grid.Fill(gruid.Cell{Rune: ' '})
		m.frame.Fill(gruid.Cell{Rune: ' '})
		return tick(m.interval)
	case timeMsg:
		m.ui.SetText(fmt.Sprintf("Pause: %t \nSpeed: %d", m.pause, m.opts.speed))
		if m.pause {
			break
		}
		if !m.pause {
			for p := range m.entities {
				m.CheckLife(p)
				around := gruid.NewRange(p.X-1, p.Y-1, p.X+2, p.Y+2)
				for p2 := range around.Points() {
					c := m.frame.At(p2)
					if c.Rune == ' ' {
						m.AddLife(p2)
					}
				}
			}
		}
		return tick(m.interval + time.Millisecond*time.Duration(m.opts.speed))
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

func (m *model) AddEntity(p gruid.Point) {
	m.entities[p] = true
}

func (m *model) RemoveEntity(p gruid.Point) {
	delete(m.entities, p)
}

func (m *model) CheckLife(p gruid.Point) {
	lifecounter := m.CountNeighbors(p)
	if lifecounter < 2 || lifecounter > 3 {
		m.RemoveEntity(p)
	} else if m.frame.At(p).Rune == ' ' {
		m.RemoveEntity(p)
	}
}

func (m *model) AddLife(p gruid.Point) {
	lifecounter := m.CountNeighbors(p)
	if lifecounter == 3 {
		m.AddEntity(p)
	}
}

func (m *model) CountNeighbors(p gruid.Point) int {
	around := gruid.NewRange(p.X-1, p.Y-1, p.X+2, p.Y+2)
	lifecounter := 0
	for p2 := range around.Points() {
		if p2 == p || !m.frame.Contains(p2) {
			continue
		} else {
			if m.frame.At(p2).Rune == '█' {
				lifecounter++
			}
		}
	}
	return lifecounter
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
	case ActionSpeedUp:
		m.opts.speed = m.opts.speed * 2
	case ActionSpeedDown:
		m.opts.speed = m.opts.speed / 2
	case ActionEnlargeMapX:
		m.opts.width++
	case ActionShrinkMapX:
		m.opts.width--
	case ActionEnlargeMapY:
		m.opts.height++
	case ActionShrinkMapY:
		m.opts.height--
	case MouseMain:
		if m.frame.At(m.action.Location).Rune == ' ' {
			m.AddEntity(m.action.Location)
		}
		m.heldAction = m.action.Type
	case MouseSecondary:
		if m.frame.At(m.action.Location).Rune == '█' {
			m.RemoveEntity(m.action.Location)
		}
		m.heldAction = m.action.Type
	case MouseRelease:
		m.heldAction = m.action.Type
	case MouseMove:
		switch m.heldAction {
		case MouseMain:
			if m.frame.At(m.action.Location).Rune == ' ' {
				m.AddEntity(m.action.Location)
			}
		case MouseSecondary:
			if m.frame.At(m.action.Location).Rune == '█' {
				m.RemoveEntity(m.action.Location)
			}
		}
	}

	switch m.action.Update {
	case UI:
		m.ui.SetText(fmt.Sprintf("Pause: %t \nSpeed: %d", m.pause, m.opts.speed))
	case Map:
	}

	return nil
}

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {
	switch msg.Key {
	case gruid.KeySpace, "p", "P":
		m.action = action{Type: ActionPause}
	case gruid.KeyEscape, "Q":
		m.action = action{Type: ActionQuit}
	case "+", "e":
		m.action = action{Type: ActionSpeedDown, Update: UI}
	case "-", "q":
		m.action = action{Type: ActionSpeedUp, Update: UI}
	case "s":
		m.action = action{Type: ActionEnlargeMapY, Update: Map}
	case "w":
		m.action = action{Type: ActionShrinkMapY, Update: Map}
	case "d":
		m.action = action{Type: ActionEnlargeMapX, Update: Map}
	case "a":
		m.action = action{Type: ActionShrinkMapX, Update: Map}
	}
}

func (m *model) updateMouse(msg gruid.MsgMouse) {
	switch msg.Action {
	case gruid.MouseMain:
		m.action = action{Type: MouseMain, Location: msg.P}
	case gruid.MouseSecondary:
		m.action = action{Type: MouseSecondary, Location: msg.P}
	case gruid.MouseRelease:
		m.action = action{Type: MouseRelease}
	case gruid.MouseMove:
		m.action = action{Type: MouseMove, Location: msg.P}
	}
}

func (m *model) Draw() gruid.Grid {
	// TODO: Go through entities here and draw them on a new frame
	newFrame := gruid.NewGrid(m.opts.width, m.opts.height)
	c := gruid.Cell{Rune: ' '}
	c.Style.Bg = ColorBackgroundSecondary
	newFrame.Fill(c)
	if len(m.entities) > 0 {
		for p := range m.entities {
			newFrame.Set(p, gruid.Cell{Rune: '█'})
		}
	}
	m.frame = newFrame
	m.grid = m.frame
	m.ui.Draw(m.grid.Slice(gruid.NewRange(0, 0, 20, 5)))
	return m.grid
}
