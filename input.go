package main

import (
	"fmt"
	"runtime"

	"codeberg.org/anaseto/gruid"
)

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

func (m *model) handleAction() gruid.Effect {

	switch m.action.Type {
	case ActionPause:
		m.pause = !m.pause
		if !m.pause {
			return tick(m.interval)
		}
		m.updateUIText(m.pause, m.opts.speed)
	case ActionQuit:
		return gruid.End()
	case ActionSpeedUp:
		m.opts.speed = m.opts.speed * 2
		m.updateUIText(m.pause, m.opts.speed)
	case ActionSpeedDown:
		if (m.opts.speed > 25) {
			m.opts.speed = m.opts.speed / 2
		}
		m.updateUIText(m.pause, m.opts.speed)
	case ActionEnlargeMapX:
		if (m.opts.width < 999) {
			m.opts.width++
		}
	case ActionShrinkMapX:
		if (m.opts.width > 30) {
			m.opts.width--
		}
	case ActionEnlargeMapY:
		if (m.opts.height < 499) {
			m.opts.height++
		}
	case ActionShrinkMapY:
		if (m.opts.height > 6) {
			m.opts.height--
		}
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

	return nil
}

func (m *model) updateUIText (pause bool, speed int) {
	m.ui.SetText(fmt.Sprintf("Pause: %t \nSpeed: %d", pause, speed))
}

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {
	switch msg.Key {
	case gruid.KeySpace, "p", "P":
		m.action = action{Type: ActionPause}
	case "+", "e":
		m.action = action{Type: ActionSpeedDown}
	case "-", "q":
		m.action = action{Type: ActionSpeedUp}
	}
	if runtime.GOOS != "js" {
		switch msg.Key {
		case gruid.KeyEscape, "Q":
			m.action = action{Type: ActionQuit}
		case "s":
			m.action = action{Type: ActionEnlargeMapY}
		case "w":
			m.action = action{Type: ActionShrinkMapY}
		case "d":
			m.action = action{Type: ActionEnlargeMapX}
		case "a":
			m.action = action{Type: ActionShrinkMapX}
		}
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

