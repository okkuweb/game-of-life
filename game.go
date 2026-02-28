package main

import (
	"fmt"
	"runtime"
	"time"

	"codeberg.org/anaseto/gruid"
)

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	m.action = action{} // reset last action information
	switch msg := msg.(type) {
	case gruid.MsgInit:
		m.frame = gruid.NewGrid(m.opts.width, m.opts.height)
		m.grid.Fill(gruid.Cell{Rune: ' ', Style: gruid.Style{Fg: ColorLife}})
		m.frame.Fill(gruid.Cell{Rune: ' ', Style: gruid.Style{Fg: ColorLife}})
		m.ui.SetText(fmt.Sprintf("Pause: %t \nSpeed: %d", m.pause, m.opts.speed))
	case timeMsg:
		if m.pause {
			break
		}
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
		return tick(m.interval + time.Millisecond*time.Duration(m.opts.speed))
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	case gruid.MsgMouse:
		m.updateMouse(msg)
	}
	return m.handleAction()
}

type timeMsg time.Time

var tickPending bool
func tick(d time.Duration) gruid.Cmd {
	if tickPending {
		return func() gruid.Msg { return nil }
	}
	tickPending = true
	t := time.NewTimer(d)
	return func() gruid.Msg {
		<-t.C
		tickPending = false
		return timeMsg{}
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

func (m *model) Draw() gruid.Grid {
	if runtime.GOOS != "js" {
		m.grid = gruid.NewGrid(1000, 500)
	} else {
		m.grid = gruid.NewGrid(m.opts.width, m.opts.height)
	}
	m.frame = gruid.NewGrid(m.opts.width, m.opts.height)
	c := gruid.Cell{Rune: ' ', Style: gruid.Style{Bg: ColorPrimary}}
	m.frame.Fill(c)
	if len(m.entities) > 0 {
		for p := range m.entities {
			m.frame.Set(p, gruid.Cell{Rune: '█', Style: gruid.Style{Fg: ColorLife}})
		}
	}
	m.grid.Copy(m.frame)
	m.ui.Draw(m.grid.Slice(gruid.NewRange(0, 0, 20, 5)))
	return m.grid
}
